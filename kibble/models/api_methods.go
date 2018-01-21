package models

import "fmt"

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

// GetGenericItem - returns a generic item
func (page Page) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     page.Title,
		Images:    page.Images,
		InnerItem: page,
	}
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
