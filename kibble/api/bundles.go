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

package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"kibble/models"
	"kibble/utils"
)

// LoadAllBundles - load all bundles
func LoadAllBundles(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	path := fmt.Sprintf("%s/services/meta/v1/bundles", cfg.SiteURL)
	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	var apiBundles []BundleV1
	err = json.Unmarshal([]byte(data), &apiBundles)
	if err != nil {
		return err
	}

	for _, b := range apiBundles {
		n := b.mapToModel(site.Config, itemIndex)
		site.Bundles = append(site.Bundles, n)
		itemIndex.Set(n.Slug, n.GetGenericItem())
	}

	return nil
}

func (b BundleV1) mapToModel(serviceConfig models.ServiceConfig, itemIndex models.ItemIndex) models.Bundle {
	return models.Bundle{
		ID:            b.ID,
		Slug:          fmt.Sprintf("/bundle/%d", b.ID),
		TitleSlug:     slug.Make(b.Title),
		Title:         b.Title,
		PromoURL:      b.PromoURL,
		PublishedDate: b.PublishedDate,
		Images: models.ImageSet{
			Portrait:   serviceConfig.ForceAbsoluteImagePath(b.PortraitImage),
			Landscape:  serviceConfig.ForceAbsoluteImagePath(b.LandscapeImage),
			Carousel:   serviceConfig.ForceAbsoluteImagePath(b.CarouselImage),
			Background: serviceConfig.ForceAbsoluteImagePath(b.BgImage),
		},
		Seo: models.Seo{
			SiteName:    serviceConfig.GetSiteName(),
			Title:       serviceConfig.GetSEOTitle(b.SeoTitle, b.Title),
			Keywords:    serviceConfig.GetKeywords(b.SeoKeywords),
			Description: utils.Coalesce(b.SeoDescription, b.Description),
			Image:       serviceConfig.SelectDefaultImageType(b.LandscapeImage, b.PortraitImage),
			VideoURL:    b.PromoURL,
		},
		Items:       itemIndex.MapToUnresolvedItems(b.Items),
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
		Description: b.Description,
		Tagline:     b.Tagline,
	}
}

// BundleV1 - model
type BundleV1 struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Tagline        string    `json:"tagline"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	PublishedDate  time.Time `json:"published_date"`
	SeoTitle       string    `json:"seo_title"`
	SeoKeywords    string    `json:"seo_keywords"`
	SeoDescription string    `json:"seo_description"`
	PortraitImage  string    `json:"portrait_image"`
	LandscapeImage string    `json:"landscape_image"`
	CarouselImage  string    `json:"carousel_image"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	BgImage        string    `json:"bg_image"`
	PromoURL       string    `json:"promo_url"`
	ExternalID     string    `json:"external_id"`
	Items          []string  `json:"items"`
}
