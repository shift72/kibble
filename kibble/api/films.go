//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"kibble/utils"

	"kibble/models"

	"github.com/gosimple/slug"
)

// loadFilmSummary - load the bios request
func loadFilmSummary(cfg *models.Config) ([]filmSummary, error) {

	var summary []filmSummary

	path := fmt.Sprintf("%s/services/meta/v2/film/index", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &summary)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// AppendAllFilms -
func AppendAllFilms(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	summary, err := loadFilmSummary(cfg)
	if err != nil {
		return err
	}

	for i := 0; i < len(summary); i++ {

		itemIndex.SetWithStatus(fmt.Sprintf("/film/%d", summary[i].ID), summary[i].StatusID, models.Unresolved)

	}

	return nil
}

// AppendFilms - load a list of films
func AppendFilms(cfg *models.Config, site *models.Site, slugs []string, itemIndex models.ItemIndex) error {
	sort.Strings(slugs)

	if len(slugs) > 300 {
		slugs = slugs[:300]
	}

	// convert /film/1,film/2 -> 1,2
	ids := strings.Replace(
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(slugs)), ","), "[]"),
		"/film/", "", -1)

	// set index to empty for the items requested
	for _, s := range slugs {
		itemIndex.Set(s, models.Empty)
	}

	path := fmt.Sprintf("%s/services/meta/v2/film/%s/show_multiple", cfg.SiteURL, ids)
	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	var details []json.RawMessage
	err = json.Unmarshal([]byte(data), &details)
	if err != nil {
		log.Error("film.error: %s", err)
		log.Debug("invalid data %s", string(data))
		return err
	}

	for i := 0; i < len(details); i++ {

		var film filmV2
		err = json.Unmarshal(details[i], &film)

		if err == nil {

			f := film.mapToModel(site.Config, itemIndex)
			site.Films = append(site.Films, f)
			itemIndex.Replace(f.Slug, f.GetGenericItem())

		} else {
			log.Error("film.error: %s", err)
			log.Debug("invalid data %s", string(details[i]))
		}
	}

	return nil
}

func (f filmV2) mapToModel(serviceConfig models.ServiceConfig, itemIndex models.ItemIndex) models.Film {

	// Convert 'foo_image' or 'foo' to 'Foo'
	for key, value := range f.ImageUrls {
		// Special case
		if (key == "bg") || (key == "bg_image") {
			key = "background"
		}

		titleCaseKey := strings.Title(strings.ToLower(key))
		titleCaseKey = strings.Replace(titleCaseKey, "_image", "", -1)

		if titleCaseKey != key {
			f.ImageUrls[titleCaseKey] = value
			delete(f.ImageUrls, key)
		}
	}

	// A couple of defaults used to fallback on
	landscapeImage := ""
	if f.ImageUrls["Landscape"] != nil {
		landscapeImage = f.ImageUrls["Landscape"].(string)
	}

	portraitImage := ""
	if f.ImageUrls["Portrait"] != nil {
		portraitImage = f.ImageUrls["Portrait"].(string)
	}

	film := models.Film{
		ID:              f.ID,
		Slug:            f.Slug,
		Title:           f.Title,
		TitleSlug:       slug.Make(f.Title),
		StatusID:        itemIndex.Get(f.Slug).StatusID,
		Overview:        f.Overview,
		Tagline:         f.Tagline,
		ReleaseDate:     utils.ParseTimeFromString(f.ReleaseDate),
		Runtime:         models.Runtime(f.Runtime),
		Countries:       f.Countries,
		Languages:       f.Languages,
		Genres:          f.Genres,
		AwardCategories: make([]models.AwardCategory, 0),
		Tags:            f.Tags,
		Seo: models.Seo{
			SiteName:    serviceConfig.GetSiteName(),
			Title:       serviceConfig.GetSEOTitle(f.SeoTitle, f.Title),
			Keywords:    serviceConfig.GetKeywords(f.SeoKeywords),
			Description: utils.Coalesce(f.SeoDescription, f.Tagline),
			Image:       serviceConfig.SelectDefaultImageType(landscapeImage, portraitImage),
		},
		Images:          f.ImageUrls,
		Recommendations: itemIndex.MapToUnresolvedItems(f.Recommendations),
		Trailers:        make([]models.Trailer, 0),
		Cast:            make([]models.CastMember, 0),
		Crew:            make([]models.CrewMember, 0),
		CustomFields:    f.CustomFields,
		Refs: models.FilmRefs{
			LetterboxdID: f.Refs.LetterboxdID,
		},
		Subtitles: f.Subtitles,
	}

	for _, s := range f.Studio {
		film.Studio = append(film.Studio, s.Name)
	}

	for key, value := range f.Classifications {
		film.Classifications = append(film.Classifications, models.Classification{
			CountryCode: key,
			Label:       value.Label,
			Description: value.Description,
		})
	}

	for _, t := range f.SubtitleTracks {
		film.SubtitleTracks = append(film.SubtitleTracks, models.SubtitleTrack{
			Language: t.Language,
			Name:     t.Name,
			Type:     t.Type,
			Path:     t.Path,
		})
	}

	// map trailers
	for i, t := range f.Trailers {
		if i == 0 {
			// set the first video URL
			film.Seo.VideoURL = t.URL
		}

		film.Trailers = append(film.Trailers, models.Trailer{
			URL:  t.URL,
			Type: t.Type,
		})
	}

	// cast
	for _, t := range f.Cast {
		film.Cast = append(film.Cast, models.CastMember{
			Name:      t.Name,
			Character: t.Character,
		})
	}

	// award categories
	for _, t := range f.AwardCategories {
		film.AwardCategories = append(film.AwardCategories, models.AwardCategory{
			Title:        t.Title,
			DisplayLabel: t.DisplayLabel,
			IsWinner:     t.IsWinner,
		})
	}

	// crew
	for _, t := range f.Crew {
		film.Crew = append(film.Crew, models.CrewMember{
			Name: t.Name,
			Job:  t.Job,
		})
	}

	// add bonuses - supports linking to bonus entries (supported??)
	for _, bonus := range f.Bonuses {
		b := bonus.mapToModel3(film.Slug, film.Images)
		film.Bonuses = append(film.Bonuses, b)
		itemIndex.Set(b.Slug, b.GetGenericItem())
	}

	// if seo image is available, use it
	seo_image := film.Images["Seo"]
	if (seo_image != nil) && (len(seo_image.(string)) > 0) {
		film.Seo.Image = seo_image.(string)
	}

	return film
}

// Film - all of the film bits
type filmV2 struct {
	Trailers []struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"trailers,omitempty"`
	Bonuses []bonusContentV2 `json:"bonuses"`
	Cast    []struct {
		Name      string `json:"name"`
		Character string `json:"character"`
	} `json:"cast"`
	Crew []struct {
		Name string `json:"name"`
		Job  string `json:"job"`
	} `json:"crew"`
	Studio []struct {
		Name string `json:"name"`
	} `json:"studio"`
	Overview        string                      `json:"overview"`
	Tagline         string                      `json:"tagline"`
	ReleaseDate     string                      `json:"release_date,omitempty"`
	Runtime         float64                     `json:"runtime"`
	Countries       []string                    `json:"countries"`
	Languages       []string                    `json:"languages"`
	Genres          []string                    `json:"genres"`
	Tags            []string                    `json:"tags"`
	Title           string                      `json:"title"`
	Slug            string                      `json:"slug"`
	FilmID          int                         `json:"film_id"`
	ID              int                         `json:"id"`
	ImageUrls       map[string]interface{}      `json:"image_urls"`
	Recommendations []string                    `json:"recommendations"`
	Subtitles       []string                    `json:"subtitles"`
	SubtitleTracks  []subtitleTrackV1           `json:"subtitle_tracks"`
	Classifications map[string]classificationV1 `json:"classifications"`
	SeoTitle        string                      `json:"seo_title"`
	SeoKeywords     string                      `json:"seo_keywords"`
	SeoDescription  string                      `json:"seo_description"`
	CustomFields    map[string]interface{}      `json:"custom"`
	Refs            struct {
		LetterboxdID string `json:"letterboxd_id"`
	} `json:"refs"`
	AwardCategories []struct {
		Title        string `json:"title"`
		DisplayLabel string `json:"display_label"`
		IsWinner     bool   `json:"is_winner"`
	} `json:"award_categories"`
}

type subtitleTrackV1 struct {
	Language string `json:"language"`
	Name     string `json:"language_name"`
	Type     string `json:"type"`
	Path     string `json:"path"`
}

type classificationV1 struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}
