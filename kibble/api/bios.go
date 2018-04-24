package api

import (
	"encoding/json"
	"fmt"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

// LoadBios - load the bios request
func LoadBios(cfg *models.Config, serviceConfig models.ServiceConfig, itemIndex models.ItemIndex) (models.Pages, models.Navigation, error) {

	var bios biosV1

	path := fmt.Sprintf("%s/services/meta/v1/bios", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return nil, bios.Navigation, err
	}

	err = json.Unmarshal([]byte(data), &bios)
	if err != nil {
		return nil, bios.Navigation, err
	}

	pages := make(models.Pages, 0)
	for _, p := range bios.Pages {
		page := p.mapToModel(serviceConfig, itemIndex)
		pages = append(pages, page)
		itemIndex.Set(page.Slug, page.GetGenericItem())
	}

	return pages, bios.Navigation, nil
}

func (p pageV1) mapToModel(serviceConfig models.ServiceConfig, itemIndex models.ItemIndex) models.Page {

	page := models.Page{
		ID:        p.ID,
		Slug:      fmt.Sprintf("/page/%d", p.ID),
		TitleSlug: p.Slug, // pages already has a property called slug which is really a title slug
		Title:     p.Title,
		Tagline:   p.Tagline,
    Content: p.Content,
		PageType:  p.PageType,
		URL:       p.URL,
		// Page images are currently relatively pathed in the bios response
		Images: models.ImageSet{
			Portrait:   serviceConfig.ForceAbsoluteImagePath(p.PortraitImage),
			Landscape:  serviceConfig.ForceAbsoluteImagePath(p.LandscapeImage),
			Carousel:   serviceConfig.ForceAbsoluteImagePath(p.CarouselImage),
			Background: serviceConfig.ForceAbsoluteImagePath(p.HeaderImage),
			Header: 		serviceConfig.ForceAbsoluteImagePath(p.HeaderImage),
		},
		PageCollections: make([]models.PageCollection, 0),
	}

	page.Seo = models.Seo{
		SiteName: serviceConfig.GetSiteName(),
		Title:    serviceConfig.GetSEOTitle(p.SeoTitle, page.Title),
		Keywords: serviceConfig.GetKeywords(p.SeoKeywords),
		Image:    serviceConfig.SelectDefaultImageType(page.Images.Landscape, page.Images.Portrait),
	}

	for _, pf := range p.PageFeatures {
		page.PageCollections = append(page.PageCollections, pf.mapToModel(serviceConfig, itemIndex))
	}

	return page
}

func (pf pageFeatureV1) mapToModel(serviceConfig models.ServiceConfig, itemIndex models.ItemIndex) models.PageCollection {
	return models.PageCollection{
		ID:          pf.FeatureID,
		Layout:      pf.Layout,
		ItemsPerRow: pf.ItemsPerRow,
		ItemLayout:  pf.ItemLayout,
		Slug:        fmt.Sprintf("/page-feature/%d", pf.FeatureID),
		TitleSlug:   pf.Slug,
		DisplayName: pf.DisplayName,
		Items:       itemIndex.MapToUnresolvedItems(pf.Items),
	}
}

type biosV1 struct {
	Navigation models.Navigation `json:"navigation"`
	Pages      []pageV1          `json:"pages"`
}

type pageFeatureV1 struct {
	FeatureID   int      `json:"feature_id"`
	Layout      string   `json:"layout"`
	ItemsPerRow int      `json:"items_per_row"`
	ItemLayout  string   `json:"item_layout"`
	Slug        string   `json:"slug"`
	DisplayName string   `json:"display_name"`
	Items       []string `json:"items"`
}

type pageV1 struct {
	CarouselImage  string          `json:"carousel_image"`
	Content        string          `json:"content"`
	HeaderImage    string          `json:"header_image"`
	ID             int             `json:"id"`
	LandscapeImage string          `json:"landscape_image"`
	PageFeatures   []pageFeatureV1 `json:"page_features"`
	PageType       string          `json:"page_type"`
	PortraitImage  string          `json:"portrait_image"`
	SeoDescription string          `json:"seo_description"`
	SeoKeywords    string          `json:"seo_keywords"`
	SeoTitle       string          `json:"seo_title"`
	Slug           string          `json:"slug"`
	Tagline        string          `json:"tagline"`
	Title          string          `json:"title"`
	URL            string          `json:"url"`
}

type filmSummary struct {
	BackgroundImage     interface{} `json:"background_image"`
	CarouselImage       interface{} `json:"carousel_image"`
	ClassificationImage interface{} `json:"classification_image"`
	HeaderImage         interface{} `json:"header_image"`
	ID                  int         `json:"id"`
	ImdbID              interface{} `json:"imdb_id"`
	LandscapeImage      interface{} `json:"landscape_image"`
	PortraitImage       string      `json:"portrait_image"`
	PublishedDate       string      `json:"published_date"`
	Slug                string      `json:"slug"`
	StatusID            int         `json:"status_id"`
	Title               string      `json:"title"`
}
