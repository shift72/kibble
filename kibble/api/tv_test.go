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
