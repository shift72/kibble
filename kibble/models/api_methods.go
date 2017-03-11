package models

var (
	// Empty - Generic Item
	Empty = GenericItem{Slug: "empty"}
	// Unresolved - a slug to an item that has not been requested
	Unresolved = GenericItem{Slug: "unresolved"}
)

// FindPageByID - find the page by id
func (pages PageCollection) FindPageByID(pageID int) (*Page, error) {
	for _, p := range pages {
		if p.ID == pageID {
			return &p, nil
		}
	}
	return nil, nil
}

// FindPageBySlug - find the page by the slug
func (pages PageCollection) FindPageBySlug(slug string) (*Page, error) {
	for _, p := range pages {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, nil
}

// FindFilmByID - find film by id
func (films FilmCollection) FindFilmByID(filmID int) (*Film, error) {
	for _, p := range films {
		if p.ID == filmID {
			return &p, nil
		}
	}
	return nil, nil
}

// FindFilmBySlug - find the film by the slug
func (films FilmCollection) FindFilmBySlug(slug string) (*Film, error) {
	for _, p := range films {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, nil
		}
	}
	return nil, nil
}

// GetGenericItem - returns a generic item
func (page Page) GetGenericItem() GenericItem {
	return GenericItem{
		Title: page.Title,
		Images: ImageSet{
			CarouselImage:   page.CarouselImage,
			PortraitImage:   page.PortraitImage,
			LandscapeImage:  page.LandscapeImage,
			BackgroundImage: page.HeaderImage,
		},
		InnerItem: page,
	}
}

// GetGenericItem - returns a generic item
func (film Film) GetGenericItem() GenericItem {
	return GenericItem{
		Title: film.Title,
		Images: ImageSet{
			CarouselImage:   &film.ImageUrls.Carousel,
			PortraitImage:   &film.ImageUrls.Portrait,
			LandscapeImage:  &film.ImageUrls.Landscape,
			BackgroundImage: &film.ImageUrls.Bg,
		},
		InnerItem: film,
	}
}

// GetGenericItem - returns a generic item
func (bundle Bundle) GetGenericItem() GenericItem {
	return GenericItem{
		Title: bundle.Title,
		Images: ImageSet{
			CarouselImage:   &bundle.CarouselImage,
			PortraitImage:   &bundle.PortraitImage,
			LandscapeImage:  &bundle.LandscapeImage,
			BackgroundImage: &bundle.BgImage,
		},
		InnerItem: bundle,
	}
}
