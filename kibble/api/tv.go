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
	"strconv"
	"strings"
	"time"

	"kibble/models"
	"kibble/utils"

	"github.com/CloudyKit/jet"
	"github.com/gosimple/slug"
)

// loadAllTVShows - loads all tv shows
func loadAllTVShows(cfg *models.Config) ([]tvShowSummaryV3, error) {

	// load all seasons
	var summary []tvShowSummaryV3

	path := fmt.Sprintf("%s/services/meta/v3/tv/index", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &summary)
	if err != nil {
		return nil, err
	}
	summary = []tvShowSummaryV3{}
	return summary, nil
}

// AppendAllTVShows -
func AppendAllTVShows(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	summary, err := loadAllTVShows(cfg)
	if err != nil {
		return err
	}

	for ti := 0; ti < len(summary); ti++ {
		fmt.Println("errrrrrrrrrrrrrrrrrrrrrrrr", summary)
		// add tv shows
		site.TVShows = append(site.TVShows, summary[ti].mapToModel())

		for si := 0; si < len(summary[ti].Seasons); si++ {
			itemIndex.Set(summary[ti].Seasons[si].Slug, models.Unresolved)
		}
	}

	return nil
}

// AppendTVSeasons - load a list of tv seasons
func AppendTVSeasons(cfg *models.Config, site *models.Site, slugs []string, itemIndex models.ItemIndex, shortCodeTmplSet *jet.Set) error {

	sort.Strings(slugs)
	ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(slugs)), ","), "[]")

	// set index to empty for the items requested
	for _, s := range slugs {
		itemIndex.Set(s, models.Empty)
	}

	path := fmt.Sprintf("%s/services/meta/v2/tv/season/show_multiple?items=%s", cfg.SiteURL, ids)
	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	var details tvSeasonShowMultipleResponseV2
	err = json.Unmarshal([]byte(data), &details)
	if err != nil {
		log.Error("tv.error: %s", err)
		log.Debug("invalid data %s", string(data))
		return err
	}

	for i := 0; i < len(details.Seasons); i++ {
		var seasonV2 tvSeasonV2
		err = json.Unmarshal(details.Seasons[i], &seasonV2)
		if err == nil {
			season := seasonV2.mapToModel(site.Config, itemIndex, shortCodeTmplSet)

			// merge the tv show information collected by the tvShowSummary and the season meta object
			// before this stage the season.ShowInfo does not have an ID field (its not returned from the api)
			if tvShowID, ok := utils.ParseIntFromSlug(season.Slug, 2); ok {
				if show, ok := site.TVShows.FindTVShowByID(tvShowID); ok {
					season.ShowInfo = mergeTVShow(show, season.ShowInfo)
					itemIndex.Set(season.ShowInfo.Slug, season.ShowInfo.GetGenericItem())
				}
			}

			// add any episodes into the global episode list
			// because the api does not return the show id, we need to re add the season here as the references somehow change
			for _, episode := range season.Episodes {
				episode.Season = &season
				site.TVEpisodes = append(site.TVEpisodes, episode)
			}

			site.TVSeasons = append(site.TVSeasons, &season)
			itemIndex.Set(season.Slug, season.GetGenericItem())
		} else {
			log.Info("Failed marshalling season %s %s", details.Seasons[i], err)
		}
	}

	// populate the tvshow with loaded seasons
	for i := 0; i < len(site.TVShows); i++ {
		for _, seasonSlug := range site.TVShows[i].AvailableSeasons {
			if season, ok := site.TVSeasons.FindTVSeasonBySlug(seasonSlug); ok {
				site.TVShows[i].Seasons = append(site.TVShows[i].Seasons, season)
			}
		}
	}

	return nil
}

func (t tvShowSummaryV3) mapToModel() *models.TVShow {
	return &models.TVShow{
		ID:        t.ID,
		Slug:      t.Slug,
		Title:     t.Title,
		TitleSlug: slug.Make(t.Title),
		Images: models.ImageSet{
			Portrait:  t.ImageUrls.Portrait,
			Landscape: t.ImageUrls.Landscape,
		},
	}
}

func mergeTVShow(tvShowA *models.TVShow, tvShowB *models.TVShow) *models.TVShow {
	tvShowA.Genres = tvShowB.Genres
	tvShowA.Overview = tvShowB.Overview
	tvShowA.Countries = tvShowB.Countries
	tvShowA.Languages = tvShowB.Languages
	tvShowA.ReleaseDate = tvShowB.ReleaseDate
	tvShowA.Studio = tvShowB.Studio
	tvShowA.AvailableSeasons = tvShowB.AvailableSeasons
	return tvShowA
}

func (t tvSeasonV2) mapToModel(serviceConfig models.ServiceConfig, itemIndex models.ItemIndex, shortCodeTmplSet *jet.Set) models.TVSeason {

	seasonNumber, _ := utils.ParseIntFromSlug(t.Slug, 4)

	//TODO: missing custom SEO fields for the season

	season := models.TVSeason{
		Slug:            t.Slug,
		Overview:        t.Overview,
		Tagline:         t.Tagline,
		PublishingState: t.PublishingState,
		Title:           t.Title,
		SeasonNumber:    seasonNumber,
		Seo: models.Seo{
			SiteName:    serviceConfig.GetSiteName(),
			Title:       serviceConfig.GetSEOTitle(t.SeoTitle, t.ShowInfo.Title),
			Keywords:    serviceConfig.GetKeywords(t.SeoKeywords),
			Description: utils.Coalesce(t.SeoDescription, t.Tagline),
			Image:       serviceConfig.SelectDefaultImageType(t.ImageUrls.Landscape, t.ImageUrls.Portrait),
		},
		Images: models.ImageSet{
			Portrait:       t.ImageUrls.Portrait,
			Landscape:      t.ImageUrls.Landscape,
			Header:         t.ImageUrls.Header,
			Carousel:       t.ImageUrls.Carousel,
			Background:     t.ImageUrls.Bg,
			Classification: t.ImageUrls.Classification,
			Seo:            t.ImageUrls.Seo,
		},
		ShowInfo:        t.ShowInfo.mapToModel(),
		Episodes:        make([]*models.TVEpisode, 0),
		Recommendations: itemIndex.MapToUnresolvedItems(t.Recommendations),
		Trailers:        make([]models.Trailer, 0),
		Cast:            make([]models.CastMember, 0),
		Crew:            make([]models.CrewMember, 0),
		CustomFields:    t.CustomFields,
	}

	// map trailers
	for i, t := range t.Trailers {
		if i == 0 {
			// set the first video URL
			season.Seo.VideoURL = t.URL
		}

		season.Trailers = append(season.Trailers, models.Trailer{
			URL:  t.URL,
			Type: t.Type,
		})
	}

	// cast
	for _, t := range t.Cast {
		season.Cast = append(season.Cast, models.CastMember{
			Name:      t.Name,
			Character: t.Character,
		})
	}

	// crew
	for _, t := range t.Crew {
		season.Crew = append(season.Crew, models.CrewMember{
			Name: t.Name,
			Job:  t.Job,
		})
	}

	// episodes
	for _, t := range t.Episodes {
		e := t.mapToModel(&season)
		season.Episodes = append(season.Episodes, &e)
		itemIndex.Set(e.Slug, e.GetGenericItem())
	}

	// add bonuses - supports linking to bonus entries (supported??)
	for _, bonus := range t.Bonuses {
		b := bonus.mapToModel2(season.Slug, season.Images, shortCodeTmplSet)
		season.Bonuses = append(season.Bonuses, b)
		itemIndex.Set(b.Slug, b.GetGenericItem())
	}

	// if seo image is available, use it
	if len(season.Images.Seo) > 0 {
		season.Seo.Image = season.Images.Seo
	}

	// classifications
	for key, value := range t.Classifications {
		season.Classifications = append(season.Classifications, models.Classification{
			CountryCode: key,
			Label:       value.Label,
			Description: value.Description,
		})
	}

	return season
}

func (t tvShowV2) mapToModel() *models.TVShow {

	show := models.TVShow{
		Title:            t.Title,
		TitleSlug:        slug.Make(t.Title),
		Genres:           t.Genres,
		Overview:         t.Overview,
		Countries:        t.Countries,
		Languages:        t.Languages,
		ReleaseDate:      utils.ParseTimeFromString(t.ReleaseDate),
		Tagline:          t.Tagline,
		Studio:           make([]string, 0),
		AvailableSeasons: t.AvailableSeasons,
	}

	for _, t := range t.Studio {
		show.Studio = append(show.Studio, t.Name)
	}

	return &show
}

func (t tvEpisodeV2) mapToModel(season *models.TVSeason) models.TVEpisode {

	episode := models.TVEpisode{
		Slug:          season.Slug + "/episode/" + strconv.Itoa(t.EpisodeNumber),
		Title:         t.Title,
		TitleSlug:     slug.Make(t.Title),
		EpisodeNumber: t.EpisodeNumber,
		Overview:      t.Overview,
		Runtime:       models.Runtime(t.Runtime),
		Images: models.ImageSet{
			Portrait:       utils.Coalesce(t.ImageUrls.Portrait, season.Images.Portrait),
			Landscape:      utils.Coalesce(t.ImageUrls.Landscape, season.Images.Landscape),
			Header:         utils.Coalesce(t.ImageUrls.Header, season.Images.Header),
			Carousel:       utils.Coalesce(t.ImageUrls.Carousel, season.Images.Carousel),
			Background:     utils.Coalesce(t.ImageUrls.Bg, season.Images.Background),
			Classification: utils.Coalesce(t.ImageUrls.Classification, season.Images.Classification),
		},
		SubtitleTracks: make([]models.SubtitleTrack, 0),
		CustomFields:   t.CustomFields,
		Season:         season,
	}

	for _, st := range t.SubtitleTracks {
		episode.SubtitleTracks = append(episode.SubtitleTracks, models.SubtitleTrack{
			Language: st.Language,
			Name:     st.Name,
			Type:     st.Type,
			Path:     st.Path,
		})
	}

	return episode
}

type tvEpisodeV2 struct {
	Title          string  `json:"title"`
	EpisodeNumber  int     `json:"episode_number"`
	DisplayTitle   string  `json:"displayTitle"`
	Overview       string  `json:"overview"`
	Runtime        float32 `json:"runtime"`
	LandscapeImage string  `json:"landscape_image"`
	ImageUrls      struct {
		Portrait       string `json:"portrait"`
		Landscape      string `json:"landscape"`
		Header         string `json:"header"`
		Carousel       string `json:"carousel"`
		Bg             string `json:"bg"`
		Classification string `json:"classification"`
	} `json:"image_urls"`
	SubtitleTracks []subtitleTrackV1      `json:"subtitle_tracks"`
	CustomFields   map[string]interface{} `json:"custom"`
}

type tvShowV2 struct {
	Trailers    []interface{} `json:"trailers"`
	Genres      []string      `json:"genres"`
	Overview    string        `json:"overview"`
	Countries   []string      `json:"countries"`
	Languages   []string      `json:"languages"`
	ReleaseDate string        `json:"release_date,omitempty"`
	Tagline     string        `json:"tagline"`
	Subtitles   string        `json:"subtitles"`
	Studio      []struct {
		Name string `json:"name"`
	} `json:"studio"`
	Title            string   `json:"title"`
	AvailableSeasons []string `json:"available_seasons"`
}

type tvSeasonV2 struct {
	Trailers []struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"trailers"`
	Episodes []tvEpisodeV2    `json:"episodes"`
	Bonuses  []bonusContentV2 `json:"bonuses"`
	Cast     []struct {
		Name      string `json:"name"`
		Character string `json:"character"`
	} `json:"cast"`
	Crew []struct {
		Name string `json:"name"`
		Job  string `json:"job"`
	} `json:"crew"`
	Tagline         string   `json:"tagline"`
	Overview        string   `json:"overview"`
	PublishingState string   `json:"publishing_state"`
	Title           string   `json:"title"`
	ShowInfo        tvShowV2 `json:"show_info"`
	Slug            string   `json:"slug"`
	SeasonNum       int      `json:"season_num"`
	ImageUrls       struct {
		Portrait       string `json:"portrait"`
		Landscape      string `json:"landscape"`
		Header         string `json:"header"`
		Carousel       string `json:"carousel"`
		Bg             string `json:"bg"`
		Classification string `json:"classification"`
		Seo            string `json:"seo"`
	} `json:"image_urls"`
	Recommendations []string                    `json:"recommendations"`
	SeoTitle        string                      `json:"seo_title"`
	SeoKeywords     string                      `json:"seo_keywords"`
	SeoDescription  string                      `json:"seo_description"`
	CustomFields    map[string]interface{}      `json:"custom"`
	Classifications map[string]classificationV1 `json:"classifications"`
}

type tvSeasonShowMultipleResponseV2 struct {
	Seasons []json.RawMessage `json:"seasons"`
}

type tvShowSummaryV3 struct {
	ID            int       `json:"id"`
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	StatusID      int       `json:"status_id"`
	PublishedDate time.Time `json:"published_date"`
	Seasons       []struct {
		Slug     string `json:"slug"`
		StatusID int    `json:"status_id"`
	} `json:"seasons"`
	ImageUrls struct {
		Portrait  string `json:"portrait_image"`
		Landscape string `json:"landscape"`
	}
}
