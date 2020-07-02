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
	"time"

	"kibble/utils"
)

// Film - all of the film bits
type Film struct {
	ID              int
	Slug            string
	Title           string
	TitleSlug       string
	Trailers        []Trailer
	Bonuses         BonusContentCollection
	Cast            []CastMember
	Crew            []CrewMember
	Studio          []string
	Overview        string
	Tagline         string
	ReleaseDate     time.Time
	Runtime         Runtime
	Countries       StringCollection
	Languages       StringCollection
	Genres          StringCollection
	Seo             Seo
	Images          ImageSet
	Recommendations []GenericItem
	Subtitles       []string
	SubtitleTracks  []SubtitleTrack
	CustomFields    CustomFields
	Classifications []Classification
}

// FilmCollection - all films
type FilmCollection []Film

// FindFilmByID - find film by id
func (films FilmCollection) FindFilmByID(filmID int) (*Film, bool) {
	for _, p := range films {
		if p.ID == filmID {
			return &p, true
		}
	}
	return nil, false
}

// FindFilmBySlug - find the film by the slug
func (films FilmCollection) FindFilmBySlug(slug string) (*Film, bool) {
	for _, p := range films {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, true
		}
	}
	return nil, false
}

// GetGenericItem - returns a generic item
func (film Film) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     film.Title,
		Slug:      film.Slug,
		Images:    film.Images,
		ItemType:  "film",
		InnerItem: film,
	}
}

// GetSubtitles - translate the Subtitles list and SubtitleTracks list into a StringCollection
func (film Film) GetSubtitles() StringCollection {
	var result StringCollection
	for _, s := range film.Subtitles {
		result = utils.AppendUnique(s, result)
	}
	for _, s := range film.SubtitleTracks {
		result = utils.AppendUnique(s.Name, result)
	}
	return result
}
