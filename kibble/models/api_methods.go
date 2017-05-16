package models

// FindPageByID - find the page by id
func (pages Pages) FindPageByID(pageID int) (*Page, error) {
	for _, p := range pages {
		if p.ID == pageID {
			return &p, nil
		}
	}
	return nil, ErrDataSourceMissing
}

// FindPageBySlug - find the page by the slug
func (pages Pages) FindPageBySlug(slug string) (*Page, error) {
	for _, p := range pages {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, ErrDataSourceMissing
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
func (films FilmCollection) FindFilmByID(filmID int) (*Film, error) {
	for _, p := range films {
		if p.ID == filmID {
			return &p, nil
		}
	}
	return nil, ErrDataSourceMissing
}

// FindFilmBySlug - find the film by the slug
func (films FilmCollection) FindFilmBySlug(slug string) (*Film, error) {
	for _, p := range films {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, nil
		}
	}
	return nil, ErrDataSourceMissing
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
