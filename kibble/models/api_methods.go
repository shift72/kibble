package models

import (
	"fmt"
	"strings"
)

var (
	// Empty - Generic Item
	Empty = GenericItem{}
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

// Add - an item
func (itemIndex ItemIndex) Add(slug string, item GenericItem) {

	slugParts := strings.Split(slug, "/")
	slugType := slugParts[1]

	// item[slugType]
	t, ok := itemIndex[slugType]
	if !ok {
		itemIndex[slugType] = make(map[string]GenericItem)
		t = itemIndex[slugType]
	}

	if _, ok := t[slug]; !ok {
		t[slug] = item
	}
}

// Get - get the slug
func (itemIndex ItemIndex) Get(slug string) (item GenericItem) {

	slugParts := strings.Split(slug, "/")
	slugType := slugParts[1]

	t, ok := itemIndex[slugType]
	if ok {
		return t[slug]
	}

	return Empty
}

// IndexItems - add an index
func (site *Site) IndexItems(itemIndex ItemIndex) {

	for _, f := range site.Films {
		for _, slug := range f.Recommendations {
			itemIndex.Add(slug, Empty)
		}
	}

	for _, p := range site.Pages {
		for _, pf := range p.PageFeatures {
			for _, slug := range pf.Items {
				itemIndex.Add(slug, Empty)
			}
		}
	}
}

// LinkItems - link the items to the specific parts
func (site *Site) LinkItems(itemIndex ItemIndex) {

	for i := range site.Films {
		for _, slug := range site.Films[i].Recommendations {
			t := itemIndex.Get(slug)
			if t != Empty {
				site.Films[i].ResolvedRecommendations = append(site.Films[i].ResolvedRecommendations, t)
			}
		}
	}

	for i := range site.Pages {
		for j := range site.Pages[i].PageFeatures {
			for _, slug := range site.Pages[i].PageFeatures[j].Items {
				t := itemIndex.Get(slug)
				if t != Empty {
					site.Pages[i].PageFeatures[j].ResolvedItems = append(site.Pages[i].PageFeatures[j].ResolvedItems, t)
				}
			}
		}
	}
}

// Print - print the item index
func (itemIndex ItemIndex) Print() {
	for t, val := range itemIndex {
		fmt.Printf("type: %s\n", t)
		for k, v := range val {
			if v == Empty {
				fmt.Printf("%s - %s : nil\n", t, k)
			} else {
				fmt.Printf("%s - %s : set\n", t, k)
			}
		}
	}
}
