package models

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

// GetGenericItem - returns a generic item based on the film bonus
func (bonus FilmBonus) GetGenericItem() GenericItem {
	return GenericItem{
		Title: bonus.Title,
		Images: ImageSet{
			CarouselImage:   &bonus.ImageUrls.Carousel,
			PortraitImage:   &bonus.ImageUrls.Portrait,
			LandscapeImage:  &bonus.ImageUrls.Landscape,
			BackgroundImage: &bonus.ImageUrls.Bg,
		},
		InnerItem: bonus,
	}
}

// FindBundleByID - find the page by id
func (bundles BundleCollection) FindBundleByID(bundleID int) (*Bundle, error) {
	for _, b := range bundles {
		if b.ID == bundleID {
			return &b, nil
		}
	}
	return nil, nil
}

// FindBundleBySlug - find the bundle by the slug
func (bundles BundleCollection) FindBundleBySlug(slug string) (*Bundle, error) {
	for _, p := range bundles {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, nil
		}
	}
	return nil, nil
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