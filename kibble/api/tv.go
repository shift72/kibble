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
	"time"

	"github.com/gosimple/slug"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/utils"
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

	return summary, nil
}

// AppendAllTVShows -
func AppendAllTVShows(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	summary, err := loadAllTVShows(cfg)
	if err != nil {
		return err
	}

	for ti := 0; ti < len(summary); ti++ {

		// add tv showws
		site.TVShows = append(site.TVShows, summary[ti].mapToModel())

		for si := 0; si < len(summary[ti].Seasons); si++ {
			itemIndex.Set(summary[ti].Seasons[si].Slug, models.Unresolved)
		}
	}

	return nil
}

// AppendTVSeasons - load a list of tv seasons
func AppendTVSeasons(cfg *models.Config, site *models.Site, slugs []string, itemIndex models.ItemIndex) error {

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
			season := seasonV2.mapToModel(site.Config, itemIndex)

			// merge the tv show information collected by the tvShowSummary and the season meta object
			if tvShowID, ok := utils.ParseIntFromSlug(season.Slug, 2); ok {
				if show, ok := site.TVShows.FindTVShowByID(tvShowID); ok {
					season.ShowInfo = mergeTVShow(show, season.ShowInfo)

					itemIndex.Set(season.ShowInfo.Slug, season.ShowInfo.GetGenericItem())
				}
			}

			site.TVSeasons = append(site.TVSeasons, season)
			itemIndex.Set(season.Slug, season.GetGenericItem())
		}
	}

	// populate the tvshow with loaded seasons
	for i := 0; i < len(site.TVShows); i++ {
		for _, seasonSlug := range site.TVShows[i].AvailableSeasons {
			if season, ok := site.TVSeasons.FindTVSeasonBySlug(seasonSlug); ok {
				site.TVShows[i].Seasons = append(site.TVShows[i].Seasons, *season)
			}
		}
	}

	return nil
}

func (t tvShowSummaryV3) mapToModel() models.TVShow {
	return models.TVShow{
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

func (t tvSeasonV2) mapToModel(serviceConfig models.ServiceConfig, itemIndex models.ItemIndex) models.TVSeason {

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
			Title:       serviceConfig.GetSEOTitle(t.Title, fmt.Sprintf("Season %d", seasonNumber)),
			Keywords:    "",
			Description: t.Tagline,
			Image:       serviceConfig.SelectDefaultImageType(t.ImageUrls.Landscape, t.ImageUrls.Portrait),
		},
		Images: models.ImageSet{
			Portrait:       t.ImageUrls.Portrait,
			Landscape:      t.ImageUrls.Landscape,
			Header:         t.ImageUrls.Header,
			Carousel:       t.ImageUrls.Carousel,
			Background:     t.ImageUrls.Bg,
			Classification: t.ImageUrls.Classification,
		},
		ShowInfo:        t.ShowInfo.mapToModel(),
		Episodes:        make([]models.TVEpisode, 0),
		Recommendations: itemIndex.MapToUnresolvedItems(t.Recommendations),
		Trailers:        make([]models.Trailer, 0),
		Cast:            make([]models.CastMember, 0),
		Crew:            make([]models.CrewMember, 0),
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
		season.Episodes = append(season.Episodes, t.mapToModel())
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
		ReleaseDate:      t.ReleaseDate,
		Tagline:          t.Tagline,
		Studio:           make([]string, 0),
		AvailableSeasons: t.AvailableSeasons,
	}

	for _, t := range t.Studio {
		show.Studio = append(show.Studio, t.Name)
	}

	return &show
}

func (t tvEpisodeV2) mapToModel() models.TVEpisode {

	episode := models.TVEpisode{
		Title:         t.Title,
		EpisodeNumber: t.EpisodeNumber,
		Overview:      t.Overview,
		Runtime:       models.Runtime(t.Runtime),
		Images: models.ImageSet{
			Portrait:       t.ImageUrls.Portrait,
			Landscape:      t.ImageUrls.Landscape,
			Header:         t.ImageUrls.Header,
			Carousel:       t.ImageUrls.Carousel,
			Background:     t.ImageUrls.Bg,
			Classification: t.ImageUrls.Classification,
		},
		Subtitles: make([]models.SubtitleTrack, 0),
	}

	for _, st := range t.SubtitleTracks {
		episode.Subtitles = append(episode.Subtitles, models.SubtitleTrack{
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
	SubtitleTracks []struct {
		Language string `json:"language"`
		Name     string `json:"language_name"`
		Type     string `json:"type"`
		Path     string `json:"path"`
	} `json:"subtitle_tracks"`
}

type tvShowV2 struct {
	Trailers    []interface{} `json:"trailers"`
	Genres      []string      `json:"genres"`
	Overview    string        `json:"overview"`
	Countries   []string      `json:"countries"`
	Languages   []string      `json:"languages"`
	ReleaseDate time.Time     `json:"release_date"`
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
	Episodes []tvEpisodeV2 `json:"episodes"`
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
	} `json:"image_urls"`
	Recommendations []string `json:"recommendations"`
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
