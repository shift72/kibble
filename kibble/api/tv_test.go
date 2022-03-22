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

func TestLoadAll(t *testing.T) {

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	itemIndex := make(models.ItemIndex)
	site := &models.Site{}

	if err := AppendAllTVShows(cfg, site, itemIndex); err != nil {
		t.Error(err)
	}

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

	assert.True(t, len(itemIndex) > 0, "itemIndex")
	//ensure episodes are added to the correct index properly
	assert.True(t, len(site.TVEpisodes) > 0, "site.TVEpisodes")
	assert.NotNil(t, site.TVEpisodes[0].Season, "First episodes Season")
	assert.NotNil(t, site.TVEpisodes[0].Season.ShowInfo, "First episodes seasons ShowInfo")
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

func TestSeasonToSeoMapWithSeoImage(t *testing.T) {

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

	imageURL := "seo_image.jpeg"
	apiSeason.ImageUrls.Seo = imageURL

	model := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, model.Seo.Image, imageURL, "should be equal")
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
			Seo            string `json:"seo"`
		}{
			Portrait:       "season-portrait.jpeg",
			Landscape:      "season-landscape.jpeg",
			Header:         "season-header.jpeg",
			Carousel:       "season-carousel.jpeg",
			Bg:             "season-background.jpeg",
			Classification: "season-classification.jpeg",
			Seo:            "season-seo.jpeg",
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

func TestSeasonClassifications(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/1",
		Title: "Season One",
		Classifications: map[string]classificationV1{
			"au": {
				Label:       "Australian Label",
				Description: "Australian Description",
			},
			"nz": {
				Label:       "NZ Label",
				Description: "NZ Description",
			},
		},
	}

	model := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 2, len(model.Classifications))

	countryCodes := [2]string{model.Classifications[0].CountryCode, model.Classifications[1].CountryCode}
	assert.Contains(t, countryCodes, "au", "expect 'au' in country codes")
	assert.Contains(t, countryCodes, "nz", "expect 'nz' in country codes")
}

func TestSeasonWithoutClassifications(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/1",
		Title: "Season One",
	}

	model := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 0, len(model.Classifications))
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

func TestEpisodesAreAddedToItemIndex(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/1",
		Title: "Season One",
		Episodes: []tvEpisodeV2{{
			EpisodeNumber: 1,
			Title:         "First Episode",
		}, {
			EpisodeNumber: 2,
			Title:         "Twoth Episode",
		}},
	}

	apiSeason.mapToModel(serviceConfig, itemIndex)
	first := itemIndex.Get("/tv/1/season/1/episode/1")
	assert.NotEqual(t, models.Empty, first)
	assert.Equal(t, "First Episode", first.Title)

	second := itemIndex.Get("/tv/1/season/1/episode/2")
	assert.NotEqual(t, models.Empty, second)
	assert.Equal(t, "Twoth Episode", second.Title)
}

func TestEpisodeHasATitleSlug(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/1",
		Title: "Season One",
		Episodes: []tvEpisodeV2{{
			EpisodeNumber: 1,
			Title:         "First Episode",
		}, {
			EpisodeNumber: 2,
			Title:         "Twoth Episode",
		}},
	}

	item := apiSeason.mapToModel(serviceConfig, itemIndex)
	assert.Equal(t, "first-episode", item.Episodes[0].TitleSlug)
	assert.Equal(t, "twoth-episode", item.Episodes[1].TitleSlug)
}

func TestEpisodeSubtitles(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/3",
		Title: "Season Three",
		Episodes: []tvEpisodeV2{{
			SubtitleTracks: []subtitleTrackV1{{
				Language: "it",
				Name:     "Italian",
				Type:     "caption",
				Path:     "/subtitles/film/49/bonus/1/it/caption-18.vtt",
			}, {
				Language: "es",
				Name:     "Spanish",
				Type:     "caption",
				Path:     "/subtitles/film/49/bonus/1/es/caption-18.vtt",
			}, {
				Language: "it",
				Name:     "Italian",
				Type:     "caption",
				Path:     "/subtitles/film/49/bonus/1/it/caption-19.vtt",
			}},
		}},
	}

	item := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 2, len(item.Episodes[0].GetSubtitles()))
	assert.Contains(t, item.Episodes[0].GetSubtitles(), "Italian")
	assert.Contains(t, item.Episodes[0].GetSubtitles(), "Spanish")
}

func TestEpisodeSubtitlesNil(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		Slug:  "/tv/1/season/3",
		Title: "Season Three",
		Episodes: []tvEpisodeV2{{
			SubtitleTracks: nil,
		}},
	}

	item := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 0, len(item.Episodes[0].GetSubtitles()))
}

func TestTvSeasonCrewJobs(t *testing.T) {

	itemIndex := make(models.ItemIndex)
	serviceConfig := commonServiceConfig()

	apiSeason := tvSeasonV2{
		SeasonNum: 2,
		Title:     "Season 2",
		Slug:      "/tv/1/season/2",
		Crew: []struct {
			Name string `json:"name"`
			Job  string `json:"job"`
		}{{
			Name: "Bryan Cranston",
			Job:  "Director",
		}, {
			Name: "Charles Haid",
			Job:  "Director",
		}, {
			Name: "Terry McDonough",
			Job:  "Director",
		}, {
			Name: "John Dahl",
			Job:  "Director",
		}, {
			Name: "Terry McDonough",
			Job:  "Director",
		}, {
			Name: "J. Roberts",
			Job:  "Writer",
		}},
	}

	season := apiSeason.mapToModel(serviceConfig, itemIndex)
	jobs := season.Crew.GetJobNames()

	assert.Equal(t, 2, len(jobs))
	assert.Contains(t, jobs, "Director")
	assert.Contains(t, jobs, "Writer")

	directors := season.Crew.GetMembers("Director")
	assert.Equal(t, 4, len(directors))
	assert.Contains(t, directors, "Bryan Cranston")
	assert.Contains(t, directors, "Terry McDonough")
	assert.Contains(t, directors, "John Dahl")
	assert.Contains(t, directors, "Charles Haid")

	cinematographers := season.Crew.GetMembers("Cinematographers")
	assert.Equal(t, 0, len(cinematographers))
}
