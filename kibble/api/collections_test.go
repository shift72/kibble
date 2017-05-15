package api

import (
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/stretchr/testify/assert"
)

func TestCollectionToSeoMap(t *testing.T) {

	itemIndex := make(models.ItemIndex)

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
	assert.Equal(t, "SHIFT72 , Collection One,  VOD", model.Seo.Title, "collection title")
	assert.Equal(t, "SHIFT72, VOD, key key key", model.Seo.Keywords, "bundle keywords")
	assert.Equal(t, "Collection description", model.Seo.Description, "collection description")
	assert.Equal(t, "portrait", model.Seo.Image, "the default seo image is portrait")
	assert.Equal(t, "", model.Seo.VideoURL, "no video is defined for the collection")

	assert.Equal(t, 2, len(itemIndex["film"]), "expect the item index to include 2 films)")
}
