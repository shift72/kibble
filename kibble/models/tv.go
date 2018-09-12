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
	Title         string
	Slug          string
	EpisodeNumber int
	Overview      string
	Runtime       Runtime
	Images        ImageSet
	Subtitles     []SubtitleTrack
	CustomFields  CustomFields
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
	Trailers        []Trailer
	Episodes        []TVEpisode
	Cast            []CastMember
	Crew            []CrewMember
	Recommendations []GenericItem
	CustomFields    CustomFields
}

// TVShowCollection -
type TVShowCollection []TVShow

// TVSeasonCollection -
type TVSeasonCollection []TVSeason

// FindTVShowByID - find tv show by id
func (shows TVShowCollection) FindTVShowByID(showID int) (*TVShow, bool) {
	for i := range shows {
		if shows[i].ID == showID {
			return &shows[i], true
		}
	}
	return nil, false
}

// FindTVShowBySlug - find the tv show by the slug
func (shows TVShowCollection) FindTVShowBySlug(slug string) (*TVShow, bool) {
	for _, p := range shows {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, true
		}
	}
	return nil, false
}

// FindTVSeasonBySlug - find the film by the slug
func (tvSeasons TVSeasonCollection) FindTVSeasonBySlug(slug string) (*TVSeason, bool) {
	for _, p := range tvSeasons {
		if p.Slug == slug {
			return &p, true
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
		ItemType:  "episode",
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
