package models

import "github.com/indiereign/shift72-kibble/kibble/utils"

// FindPageByID - find the page by id
func (pages PageCollection) FindPageByID(pageID int) (*Page, error) {
	for _, p := range pages {
		if p.ID == pageID {
			return &p, nil
		}
	}
	return nil, DataSourceMissing
}

// FindPageBySlug - find the page by the slug
func (pages PageCollection) FindPageBySlug(slug string) (*Page, error) {
	for _, p := range pages {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, DataSourceMissing
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
	return nil, DataSourceMissing
}

// FindFilmBySlug - find the film by the slug
func (films FilmCollection) FindFilmBySlug(slug string) (*Film, error) {
	for _, p := range films {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, nil
		}
	}
	return nil, DataSourceMissing
}

// GetGenericItem - returns a generic item
func (film Film) GetGenericItem() GenericItem {
	return GenericItem{
		Title: film.Title,
		Slug:  film.Slug,
		Images: ImageSet{
			CarouselImage:   film.ImageUrls.Carousel,
			PortraitImage:   film.ImageUrls.Portrait,
			LandscapeImage:  film.ImageUrls.Landscape,
			BackgroundImage: film.ImageUrls.Bg,
		},
		ItemType:  "film",
		InnerItem: film,
	}
}

// GetSEO - get the film seo
func (film Film) GetSEO(config ServiceConfig) *Seo {

	return &Seo{
		Title:    config["seo_title_prefix"] + film.Title + config["seo_title_suffix"],
		Keywords: utils.Join(", ", config["seo_site_keywords"], film.Keywords),
	}
}

// GetGenericItem - returns a generic item based on the film bonus
func (bonus FilmBonus) GetGenericItem() GenericItem {
	return GenericItem{
		Title: bonus.Title,
		// Slug:  fmt.Sprintf("/bonus/%d", bonus.Number),
		Images: ImageSet{
			CarouselImage:   bonus.ImageUrls.Carousel,
			PortraitImage:   bonus.ImageUrls.Portrait,
			LandscapeImage:  bonus.ImageUrls.Landscape,
			BackgroundImage: bonus.ImageUrls.Bg,
		},
		ItemType:  "bonus",
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
	return nil, DataSourceMissing
}

// FindBundleBySlug - find the bundle by the slug
func (bundles BundleCollection) FindBundleBySlug(slug string) (*Bundle, error) {
	for _, p := range bundles {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, nil
		}
	}
	return nil, DataSourceMissing
}

// GetGenericItem - returns a generic item
func (bundle Bundle) GetGenericItem() GenericItem {
	return GenericItem{
		Title: bundle.Title,
		Images: ImageSet{
			CarouselImage:   bundle.CarouselImage,
			PortraitImage:   bundle.PortraitImage,
			LandscapeImage:  bundle.LandscapeImage,
			BackgroundImage: bundle.BgImage,
		},
		InnerItem: bundle,
	}
}

// FindCollectionByID - find the page by id
func (collections CollectionCollection) FindCollectionByID(collectionID int) (*Collection, error) {
	for _, b := range collections {
		if b.ID == collectionID {
			return &b, nil
		}
	}
	return nil, DataSourceMissing
}

// FindCollectionBySlug - find the collection by the slug
func (collections CollectionCollection) FindCollectionBySlug(slug string) (*Collection, error) {
	for _, p := range collections {
		if p.Slug == slug || p.TitleSlug == slug {
			return &p, nil
		}
	}
	return nil, DataSourceMissing
}
