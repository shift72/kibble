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

	"github.com/indiereign/shift72-kibble/kibble/utils"
	"github.com/nicksnyder/go-i18n/i18n"
)

// String - overrides String function for []string which joins all items in the array into english readable format.
func (strings StringCollection) String() string {
	return strings.Join(", ")
}

// Join - Joins all items in []string using specified separator.
func (strings StringCollection) Join(separator string) string {
	return utils.Join(separator, strings...)
}

// Hours - returns the basic number of hours in the runtime
func (rt Runtime) Hours() int {
	return int(rt) / 60
}

// Minutes - returns the number of minutes (left =over from any hours) in the runtime.
// This is not the total minutes.
func (rt Runtime) Minutes() int {
	return int(rt) % 60
}

// Localise - returns the localised runtime format
func (rt Runtime) Localise(T i18n.TranslateFunc) string {
	arg := map[string]interface{}{
		"Hours":   rt.Hours(),
		"Minutes": rt.Minutes(),
	}

	h := rt.Hours()
	if h == 0 {
		// zero != 0 in languages
		// https://github.com/nicksnyder/go-i18n/issues/58
		return T("runtime_minutes_only", arg)
	}

	return T("runtime", h, arg)
}

// FindPageByID - find the page by id
func (pages Pages) FindPageByID(pageID int) (*Page, bool) {
	for _, p := range pages {
		if p.ID == pageID {
			return &p, true
		}
	}
	return nil, false
}

// FindPageBySlug - find the page by the slug
func (pages Pages) FindPageBySlug(slug string) (*Page, bool) {
	for _, p := range pages {
		if p.Slug == slug {
			return &p, true
		}
	}
	return nil, false
}

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
func (page Page) GetGenericItem() GenericItem {
	return GenericItem{
		Slug:      page.Slug,
		Title:     page.Title,
		Images:    page.Images,
		ItemType:  "page",
		InnerItem: page,
	}
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

// GetGenericItem - returns a generic item based on the film bonus
func (bonus FilmBonus) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     bonus.Title,
		Slug:      bonus.Slug,
		Images:    bonus.Images,
		ItemType:  "bonus",
		InnerItem: bonus,
	}
}

// GetTitle - returns the title in the current language
// expect to be called as item.GetTitle(i18n) where i18n is the translation function
// for the current language
func (i *GenericItem) GetTitle(T i18n.TranslateFunc) string {
	switch i.ItemType {
	case "tvseason":
		if s, ok := i.InnerItem.(TVSeason); ok {
			return T(i.ItemType, map[string]interface{}{
				"ShowInfo": s.ShowInfo,
				"Season":   s,
			})
		}
	}
	return i.Title
}
