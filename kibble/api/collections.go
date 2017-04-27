package api

import (
	"encoding/json"
	"fmt"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

// LoadAllCollections - load all collections
func LoadAllCollections(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	path := fmt.Sprintf("%s/services/meta/v4/featured", cfg.SiteURL)
	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	details := []CollectionV4{}
	err = json.Unmarshal([]byte(data), &details)
	if err != nil {
		return err
	}

	for _, c := range details {

		collection := models.Collection{
			ID:          c.ID,
			Slug:        fmt.Sprintf("/collection/%d", c.ID),
			Title:       c.Title,
			TitleSlug:   c.TitleSlug,
			Description: c.Description,
			DisplayName: c.DisplayName,
			ItemLayout:  c.ItemLayout,
			ItemsPerRow: c.ItemsPerRow,
			ImageSet: models.ImageSet{
				LandscapeImage: c.LandscapeImage,
				PortraitImage:  c.PortraitImage,
				CarouselImage:  c.CarouselImage,
				HeaderImage:    c.HeaderImage,
			},
			Items:          c.Items, //TODO: map to generic item
			SearchQuery:    c.SearchQuery,
			SeoDescription: c.SeoDescription,
			SeoKeywords:    c.SeoKeywords,
			SeoTitle:       c.SeoTitle,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		}

		// add items
		for _, slug := range collection.Items {
			itemIndex.Set(slug, models.Unresolved)
		}

		site.Collections = append(site.Collections, collection)
	}

	return nil
}

// CollectionV4 - mapped from the v4 api
type CollectionV4 struct {
	CarouselImage  string   `json:"carousel_image"`
	CreatedAt      string   `json:"created_at"`
	Description    string   `json:"description"`
	DisplayName    string   `json:"display_name"`
	HeaderImage    string   `json:"header_image"`
	ID             int      `json:"id"`
	ItemLayout     string   `json:"item_layout"`
	Items          []string `json:"items"`
	ItemsPerRow    int      `json:"items_per_row"`
	LandscapeImage string   `json:"landscape_image"`
	PortraitImage  string   `json:"portrait_image"`
	SearchQuery    string   `json:"search_query"`
	SeoDescription string   `json:"seo_description"`
	SeoKeywords    string   `json:"seo_keywords"`
	SeoTitle       string   `json:"seo_title"`
	Title          string   `json:"title"`
	TitleSlug      string   `json:"title_slug"`
	UpdatedAt      string   `json:"updated_at"`
}
