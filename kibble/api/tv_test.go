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
	"fmt"
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadAll(t *testing.T) {

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	itemIndex := make(models.ItemIndex)
	site := &models.Site{}

	AppendAllTVShows(cfg, site, itemIndex)

}

func TestLoadTVSeasons(t *testing.T) {

	// if testing.Short() {
	// 	return
	// }

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	itemIndex := make(models.ItemIndex)
	site := &models.Site{}
	slugs := []string{
		"/tv/4/season/2",
		"/tv/41/season/1",
		"/tv/9/season/1",
	}

	err := AppendTVSeasons(cfg, site, slugs, itemIndex)
	if err != nil {
		t.Error(err)
	}

	if len(itemIndex) == 0 {
		t.Error("Expected some values to be loaded")
	}

	fmt.Printf("here shows %d ", len(site.TVShows))

}
func TestSeasonToSeoMap(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:           "/tv/1/season/1",
		Title:          "Season One",
		Overview:       "Season overview",
		SeoTitle:       "Season Season Season",
		SeoKeywords:    "key key key",
		SeoDescription: "One Season to rule them all",
		ShowInfo: tvShowV2{
			Title: "Show One",
		},
	}

	apiSeason.ImageUrls.Portrait = "portrait"

	model := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "Film On Demand", model.Seo.SiteName, "Season site name")
	assert.Equal(t, "SHIFT72 , Season Season Season , VOD", model.Seo.Title, "Season title")
	assert.Equal(t, "SHIFT72, VOD, key key key", model.Seo.Keywords, "Season keywords")
	assert.Equal(t, "One Season to rule them all", model.Seo.Description, "Season description")
	assert.Equal(t, "portrait", model.Seo.Image, "the default seo image is portrait")
}

func TestTVEpisodeSlugPopulations(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:           "/tv/1/season/1",
		Title:          "Season One",
		Overview:       "Season overview",
		SeoTitle:       "Season Season Season",
		SeoKeywords:    "key key key",
		SeoDescription: "One Season to rule them all",
		ShowInfo: tvShowV2{
			Title: "Show One",
		},
		Episodes: []tvEpisodeV2{{
			EpisodeNumber: 1,
			Title:         "First Episode",
		}, {
			EpisodeNumber: 2,
			Title:         "Second Episode",
		}},
	}

	model := apiSeason.mapToModel(serviceConfig, itemIndex)
	assert.Equal(t, "/tv/1/season/1/episode/1", model.Episodes[0].Slug, "first episode slug")
	assert.Equal(t, "/tv/1/season/1/episode/2", model.Episodes[1].Slug, "second episode slug")

}

func TestEpisodeImageFallback(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/1",
		Title: "Season One",
		ImageUrls: struct {
			Portrait       string `json:"portrait"`
			Landscape      string `json:"landscape"`
			Header         string `json:"header"`
			Carousel       string `json:"carousel"`
			Bg             string `json:"bg"`
			Classification string `json:"classification"`
		}{
			Portrait:       "season-portrait.jpeg",
			Landscape:      "season-landscape.jpeg",
			Header:         "season-header.jpeg",
			Carousel:       "season-carousel.jpeg",
			Bg:             "season-background.jpeg",
			Classification: "season-classification.jpeg",
		},
		ShowInfo: tvShowV2{
			Title: "Show One",
		},
		Episodes: []tvEpisodeV2{{
			EpisodeNumber: 1,
			Title:         "First Episode",
		}},
	}

	model := apiSeason.mapToModel(serviceConfig, itemIndex)
	assert.Equal(t, "season-portrait.jpeg", model.Episodes[0].Images.Portrait)
	assert.Equal(t, "season-landscape.jpeg", model.Episodes[0].Images.Landscape)
	assert.Equal(t, "season-header.jpeg", model.Episodes[0].Images.Header)
	assert.Equal(t, "season-carousel.jpeg", model.Episodes[0].Images.Carousel)
	assert.Equal(t, "season-background.jpeg", model.Episodes[0].Images.Background)
	assert.Equal(t, "season-classification.jpeg", model.Episodes[0].Images.Classification)
}

func TestSeasonCustomFieldSupport(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/1",
		Title: "Season One",
		CustomFields: map[string]interface{}{
			"facebook_url": "https://www.facebook.com/custompage",
			"some_key":     1,
			"another_key":  false,
		},
	}

	model := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 1, model.CustomFields["some_key"])
	assert.Equal(t, "https://www.facebook.com/custompage", model.CustomFields["facebook_url"])
	assert.Equal(t, false, model.CustomFields["another_key"])
	assert.Equal(t, nil, model.CustomFields["where is it"])
}

func TestEpisodeCustomFieldSupport(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/1",
		Title: "Season One",
		Episodes: []tvEpisodeV2{{
			EpisodeNumber: 1,
			Title:         "First Episode",
			CustomFields: map[string]interface{}{
				"facebook_url": "https://www.facebook.com/custompage",
				"some_key":     1,
				"another_key":  false,
			},
		}},
	}

	model := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 1, model.Episodes[0].CustomFields["some_key"])
	assert.Equal(t, "https://www.facebook.com/custompage", model.Episodes[0].CustomFields["facebook_url"])
	assert.Equal(t, false, model.Episodes[0].CustomFields["another_key"])
	assert.Equal(t, nil, model.Episodes[0].CustomFields["where is it"])

	assert.Equal(t, nil, model.CustomFields["hello?"])
}

func TestBonusContentModelBinding(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:     "/tv/12/season/4",
		Title:    "Season Fourth",
		Overview: "Season overview",
		ShowInfo: tvShowV2{
			Title: "Show Twelth",
		},
		Bonuses: []bonusContentV2{{
			Number: 1,
			Title:  "Behind the scenes",
			ImageUrls: struct {
				Portrait       string `json:"portrait"`
				Landscape      string `json:"landscape"`
				Header         string `json:"header"`
				Carousel       string `json:"carousel"`
				Bg             string `json:"bg"`
				Classification string `json:"classification"`
			}{
				Portrait:       "portrait",
				Landscape:      "landscape",
				Classification: "classification",
			},
		}},
	}

	model := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 1, len(model.Bonuses), "expect 1 bonus")
	assert.Equal(t, "/tv/12/season/4/bonus/1", model.Bonuses[0].Slug, "bonus.slug")

	assert.Equal(t, "portrait", model.Bonuses[0].Images.Portrait)
	assert.Equal(t, "landscape", model.Bonuses[0].Images.Landscape)
}

