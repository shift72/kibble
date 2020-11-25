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
	"testing"

	"kibble/models"

	"github.com/stretchr/testify/assert"
)

func TestPageToSeoMap(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiPage := pageV1{
		Title:          "Page One",
		SeoTitle:       "Special Page One",
		SeoDescription: "Page One is so special it hurts",
		SeoKeywords:    "key key key",
		PortraitImage:  "/portrait",
		LandscapeImage: "/landscape",
	}

	model := apiPage.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "Film On Demand", model.Seo.SiteName, "site name")
	assert.Equal(t, "SHIFT72 , Special Page One , VOD", model.Seo.Title, "page title")
	assert.Equal(t, "SHIFT72, VOD, key key key", model.Seo.Keywords, "keywords")
	assert.Equal(t, "https://s3-bla-bla/portrait", model.Seo.Image, "the default seo image is portrait")
	assert.Equal(t, "Page One is so special it hurts", model.Seo.Description, "seo description")
}

func TestPagehasAbsoluteImagePaths(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiPage := pageV1{
		Title:          "Page One",
		SeoKeywords:    "key key key",
		PortraitImage:  "/portrait",
		LandscapeImage: "/landscape",
		CarouselImage:  "/carousel",
		HeaderImage:    "/header",
	}

	model := apiPage.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "https://s3-bla-bla/portrait", model.Images.Portrait, "portrait")
	assert.Equal(t, "https://s3-bla-bla/landscape", model.Images.Landscape, "landscape")
	assert.Equal(t, "https://s3-bla-bla/carousel", model.Images.Carousel, "carousel")
	assert.Equal(t, "https://s3-bla-bla/header", model.Images.Header, "header")
	assert.Equal(t, "https://s3-bla-bla/header", model.Images.Background, "background")
}

func TestPageToPageFeatures(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiPage := pageV1{
		ID:             123,
		Title:          "Page One",
		Slug:           "page-one",
		SeoKeywords:    "key key key",
		PortraitImage:  "portrait",
		LandscapeImage: "landscape",
		PageFeatures: []pageFeatureV1{
			pageFeatureV1{
				FeatureID:   120,
				Layout:      "slider",
				ItemsPerRow: 3,
				ItemLayout:  "landscape",
				Slug:        "/page-feature/blam",
				DisplayName: "New Releases",
				Items: []string{
					"/film/1",
					"/film/2",
					"/bundle/1",
				},
			},
		},
	}

	model := apiPage.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "Page One", model.Title, "title")
	assert.Equal(t, "page-one", model.TitleSlug, "title slug")
	assert.Equal(t, "/page/123", model.Slug, "slug")

	// page features
	assert.Equal(t, "/collection/120", model.PageCollections[0].Slug)
	assert.Equal(t, "/page-feature/blam", model.PageCollections[0].TitleSlug)

	assert.Equal(t, 2, len(itemIndex["film"]), "expect the item index to include 2 films")
	assert.Equal(t, 1, len(itemIndex["bundle"]), "expect the item index to include 1 bundles")
}
