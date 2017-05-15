package api

import (
	"encoding/json"
	"fmt"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/utils"
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
		collection := c.mapToModel(site.Config, itemIndex)
		site.Collections = append(site.Collections, collection)
		itemIndex.Set(collection.Slug, collection.GetGenericItem())
	}

	return nil
}

func (c CollectionV4) mapToModel(serviceConfig models.ServiceConfig, itemIndex models.ItemIndex) models.Collection {

	return models.Collection{
		ID:          c.ID,
		Slug:        fmt.Sprintf("/collection/%d", c.ID),
		Title:       c.Title,
		TitleSlug:   c.TitleSlug,
		Description: c.Description,
		DisplayName: c.DisplayName,
		ItemLayout:  c.ItemLayout,
		ItemsPerRow: c.ItemsPerRow,
		Images: models.ImageSet{
			Landscape: c.LandscapeImage,
			Portrait:  c.PortraitImage,
			Carousel:  c.CarouselImage,
			Header:    c.HeaderImage,
		},
		Seo: models.Seo{
			SiteName:    serviceConfig.GetSiteName(),
			Title:       serviceConfig.GetSEOTitle(c.SeoTitle, c.Title),
			Keywords:    serviceConfig.GetKeywords(c.SeoKeywords),
			Description: utils.Coalesce(c.SeoDescription, c.Description),
			Image:       serviceConfig.SelectDefaultImageType(c.LandscapeImage, c.PortraitImage),
		},
		SearchQuery: c.SearchQuery,
		Items:       itemIndex.MapToUnresolvedItems(c.Items),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
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
