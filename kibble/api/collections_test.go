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

func TestCollectionToSeoMap(t *testing.T) {

	itemIndex := models.NewItemIndex()

	serviceConfig := commonServiceConfig()

	apiCollection := CollectionV4{
		Title:          "Collection One",
		Description:    "Collection description",
		SeoKeywords:    "key key key",
		PortraitImage:  "portrait",
		LandscapeImage: "landscape",
		Items:          []string{"/film/1", "/film/2"},
	}

	model := apiCollection.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "Film On Demand", model.Seo.SiteName, "collection site name")
	assert.Equal(t, "SHIFT72 , Collection One , VOD", model.Seo.Title, "collection title")
	assert.Equal(t, "SHIFT72, VOD, key key key", model.Seo.Keywords, "bundle keywords")
	assert.Equal(t, "Collection description", model.Seo.Description, "collection description")
	assert.Equal(t, "portrait", model.Seo.Image, "the default seo image is portrait")
	assert.Equal(t, "", model.Seo.VideoURL, "no video is defined for the collection")

	assert.Equal(t, 2, len(itemIndex.Items["film"]), "expect the item index to include 2 films)")
}
