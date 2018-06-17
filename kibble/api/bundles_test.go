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

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/stretchr/testify/assert"
)

func commonServiceConfig() models.ServiceConfig {
	return models.ServiceConfig{
		"image_root_path":       "https://s3-bla-bla",
		"portrait_poster_path":  "/posters-and-backdrops/282x422",
		"landscape_poster_path": "/posters-and-backdrops/380x210",
		"seo_title_prefix":      "SHIFT72 ",
		"seo_title_suffix":      " VOD",
		"seo_site_keywords":     "SHIFT72, VOD",
		"seo_site_name":         "Film On Demand",
	}
}

func TestBundleToSeoMap(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiBundle := BundleV1{
		Title:          "Bundle One",
		Description:    "Bundle description",
		SeoKeywords:    "key key key",
		PortraitImage:  "portrait",
		LandscapeImage: "landscape",
		PromoURL:       "https://video",
	}

	model := apiBundle.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "Film On Demand", model.Seo.SiteName, "bundle site name")
	assert.Equal(t, "SHIFT72 , Bundle One,  VOD", model.Seo.Title, "bundle title")
	assert.Equal(t, "SHIFT72, VOD, key key key", model.Seo.Keywords, "bundle keywords")
	assert.Equal(t, "Bundle description", model.Seo.Description, "bundle description")
	assert.Equal(t, "portrait", model.Seo.Image, "the default seo image is portrait")
	assert.Equal(t, "https://video", model.Seo.VideoURL, "video url is mapped from the PromoURL")
}
func TestBundlesToSeoDefaultImage(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := models.ServiceConfig{
		"default_image_type": "landscape",
	}

	apiBundle := BundleV1{
		PortraitImage:  "portrait",
		LandscapeImage: "landscape",
	}

	model := apiBundle.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "landscape", model.Seo.Image, "the default seo image is landscape")
}

func TestBundlesApiToModel(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiBundle := BundleV1{
		ID:             123,
		Title:          "Bundle One",
		Description:    "Bundle description",
		SeoKeywords:    "key key key",
		PortraitImage:  "portrait",
		LandscapeImage: "landscape",
		PromoURL:       "https://video",
		Items:          []string{"/film/1", "/film/2"},
	}

	model := apiBundle.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "bundle-one", model.TitleSlug, "title slug")
	assert.Equal(t, "/bundle/123", model.Slug, "slug")

	assert.Equal(t, 2, len(model.Items), "expect 2 generic items")
	assert.Equal(t, "/film/1", model.Items[0].Slug, "expect /film/1")
	assert.Equal(t, "/film/2", model.Items[1].Slug, "expect /film/2")
	assert.Equal(t, nil, model.Items[0].InnerItem, "expect inner item to be nil")

	assert.Equal(t, 2, len(itemIndex["film"]), "expect the item index to include 2 films")
}
