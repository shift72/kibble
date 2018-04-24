package api

import (
  "testing"

  "github.com/indiereign/shift72-kibble/kibble/models"
  "github.com/stretchr/testify/assert"
)

func TestPageToSeoMap(t *testing.T) {

  itemIndex := make(models.ItemIndex)

  serviceConfig := commonServiceConfig()

  apiPage := pageV1{
    Title:          "Page One",
    SeoKeywords:    "key key key",
    PortraitImage:  "/portrait",
    LandscapeImage: "/landscape",
  }

  model := apiPage.mapToModel(serviceConfig, itemIndex)

  assert.Equal(t, "Film On Demand", model.Seo.SiteName, "site name")
  assert.Equal(t, "SHIFT72 , Page One,  VOD", model.Seo.Title, "page title")
  assert.Equal(t, "SHIFT72, VOD, key key key", model.Seo.Keywords, "keywords")
  assert.Equal(t, "https://s3-bla-bla/portrait", model.Seo.Image, "the default seo image is portrait")
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
  assert.Equal(t, "/page-feature/120", model.PageCollections[0].Slug)
  assert.Equal(t, "/page-feature/blam", model.PageCollections[0].TitleSlug)

  assert.Equal(t, 2, len(itemIndex["film"]), "expect the item index to include 2 films")
  assert.Equal(t, 1, len(itemIndex["bundle"]), "expect the item index to include 1 bundles")
}
