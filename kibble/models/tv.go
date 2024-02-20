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

package models

import (
	"fmt"
	"time"

	"kibble/utils"

	"github.com/nicksnyder/go-i18n/i18n"
)

// TVShow -
type TVShow struct {
	ID               int
	Slug             string
	Trailers         []Trailer
	Genres           StringCollection
	Overview         string
	Countries        StringCollection
	Languages        StringCollection
	ReleaseDate      time.Time
	Tagline          string
	Studio           []string
	Title            string
	TitleSlug        string
	AvailableSeasons []string           `json:"-"`
	Seasons          TVSeasonCollection `json:"-"`
	Images           ImageSet
}

// TVEpisode -
type TVEpisode struct {
	Title          string
	Slug           string
	TitleSlug      string
	EpisodeNumber  int
	Overview       string
	Runtime        Runtime
	Images         ImageSet
	SubtitleTracks []SubtitleTrack
	CustomFields   CustomFields
	Season         *TVSeason
	Available      Period
}

// TVSeason -
type TVSeason struct {
	Slug         string
	SeasonNumber int
	//TODO: consider removing this
	Title           string
	Tagline         string
	Overview        string
	PublishingState string
	ShowInfo        *TVShow
	Seo             Seo
	Images          ImageSet
	Prices          PriceInfo
	Available       Period
	Trailers        []Trailer
	Episodes        TVEpisodeCollection
	Bonuses         BonusContentCollection
	Cast            []CastMember
	Crew            CrewMembers
	Recommendations []GenericItem
	CustomFields    CustomFields
	Classifications []Classification
	CarouselFocus   string
}

// TVShowCollection -
type TVShowCollection []*TVShow

// TVSeasonCollection -
type TVSeasonCollection []*TVSeason

// TVEpisodeCollection is an array of episodes
type TVEpisodeCollection []*TVEpisode

// FindTVShowByID - find tv show by id
func (shows TVShowCollection) FindTVShowByID(showID int) (*TVShow, bool) {
	coll := shows
	for i := 0; i < len(coll); i++ {
		if coll[i].ID == showID {
			return coll[i], true
		}
	}
	return nil, false
}

// FindTVShowBySlug - find the tv show by the slug
func (shows TVShowCollection) FindTVShowBySlug(slug string) (*TVShow, bool) {
	coll := shows
	for i := 0; i < len(coll); i++ {
		if coll[i].Slug == slug || coll[i].TitleSlug == slug {
			return coll[i], true
		}
	}
	return nil, false
}

// FindTVSeasonBySlug - find the film by the slug
func (tvSeasons TVSeasonCollection) FindTVSeasonBySlug(slug string) (*TVSeason, bool) {
	coll := tvSeasons
	for i := 0; i < len(coll); i++ {
		if coll[i].Slug == slug {
			return coll[i], true
		}
	}
	return nil, false
}

// FindTVEpisodeBySlug returns an episode based on the specified slug
func (episodes TVEpisodeCollection) FindTVEpisodeBySlug(slug string) (*TVEpisode, bool) {
	coll := episodes
	for i := 0; i < len(coll); i++ {
		if coll[i].Slug == slug {
			return coll[i], true
		}
	}
	return nil, false
}

// GetGenericItem - returns a generic item
func (show TVShow) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     show.Title,
		Slug:      show.Slug,
		Images:    show.Images,
		ItemType:  "tvshow",
		InnerItem: show,
	}
}

// GetGenericItem - returns a generic item
func (season TVSeason) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     fmt.Sprintf("%s - Season - %d", season.ShowInfo.Title, season.SeasonNumber),
		Slug:      season.Slug,
		Images:    season.Images,
		ItemType:  "tvseason",
		InnerItem: season,
	}
}

// GetGenericItem returns a generic item for the specific episode
func (episode TVEpisode) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     episode.Title,
		Slug:      episode.Slug,
		Images:    episode.Images,
		ItemType:  "tvepisode",
		InnerItem: episode,
	}
}

// GetTitle return the localised version of the season title
func (season TVSeason) GetTitle(T i18n.TranslateFunc) string {
	return T("tvseason", map[string]interface{}{
		"ShowInfo": *season.ShowInfo,
		"Season":   season,
	})
}

// GetTranslatedTitle returns an i18n version of a season title using the specified key as the template
func (season TVSeason) GetTranslatedTitle(T i18n.TranslateFunc, i18nKey string) string {
	if i18nKey == "" {
		i18nKey = "tvseason"
	}

	return T(i18nKey, map[string]interface{}{
		"ShowInfo": *season.ShowInfo,
		"Season":   season,
	})
}

// GetTitle return the localised version of the episode title
func (episode TVEpisode) GetTitle(T i18n.TranslateFunc) string {
	return T("tvepisode", map[string]interface{}{
		"ShowInfo": *episode.Season.ShowInfo,
		"Season":   *episode.Season,
		"Episode":  episode,
	})
}

// GetTranslatedTitle returns an i18n version of an episode title using the specified key as the template
func (episode TVEpisode) GetTranslatedTitle(T i18n.TranslateFunc, i18nKey string) string {
	if i18nKey == "" {
		i18nKey = "tvepisode"
	}

	return T(i18nKey, map[string]interface{}{
		"ShowInfo": *episode.Season.ShowInfo,
		"Season":   *episode.Season,
		"Episode":  episode,
	})
}

// GetSubtitles - translate the SubtitleTracks list into a StringCollection
func (episode TVEpisode) GetSubtitles() StringCollection {
	var result StringCollection
	for _, s := range episode.SubtitleTracks {
		result = utils.AppendUnique(s.Name, result)
	}
	return result
}
