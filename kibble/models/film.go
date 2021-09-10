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
	"sort"
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
	Crew            CrewMembers
	Studio          []string
	Overview        string
	Tagline         string
	ReleaseDate     time.Time
	Runtime         Runtime
	Countries       StringCollection
	Languages       StringCollection
	Genres          StringCollection
	Tags            StringCollection
	Seo             Seo
	Images          ImageSet
	Prices          PriceInfo
	Available       Period
	Recommendations []GenericItem
	Subtitles       []string
	SubtitleTracks  []SubtitleTrack
	CustomFields    CustomFields
	Refs            CustomFields
	Classifications []Classification
}

// FilmCollection - all films
type FilmCollection []Film

// FindFilmByID - find film by id
func (films *FilmCollection) FindFilmByID(filmID int) (*Film, bool) {
	coll := *films
	for i := 0; i < len(coll); i++ {
		if coll[i].ID == filmID {
			return &coll[i], true
		}
	}
	return nil, false
}

// FindFilmBySlug - find the film by the slug
func (films *FilmCollection) FindFilmBySlug(slug string) (*Film, bool) {
	coll := *films
	for i := 0; i < len(coll); i++ {
		if coll[i].Slug == slug || coll[i].TitleSlug == slug {
			return &coll[i], true
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

// MakeTitleSlugsUnique scans the films looking for duplicates
func (films *FilmCollection) MakeTitleSlugsUnique() {

	groups := make(map[string][]int, 0)

	// create a grouping of slugs to films first
	for fi, film := range *films {
		if groups[film.TitleSlug] == nil {
			groups[film.TitleSlug] = []int{fi}
		} else {
			groups[film.TitleSlug] = append(groups[film.TitleSlug], fi)
		}
	}

	// if any groups are larger than 1 then make them unique
	for _, group := range groups {
		if len(group) == 1 {
			continue
		}

		// sort them by id, so the first film is not changed
		sort.Slice(group, func(i int, j int) bool {
			return (*films)[group[i]].ID < (*films)[group[j]].ID
		})

		// append i + 1 to end of slug
		for j := 0; j < len(group); j++ {
			if j == 0 {
				continue
			}

			// find the correct value to update, without making a copy
			(*films)[group[j]].TitleSlug = fmt.Sprintf("%s-%d", (*films)[group[j]].TitleSlug, j+1)
		}
	}
}
