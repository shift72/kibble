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

func TestLoadFilms(t *testing.T) {

	// if testing.Short() {
	// 	return
	// }

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	summary, err := loadFilmSummary(cfg)
	if err != nil {
		t.Error(err)
	}

	if len(summary) == 0 {
		t.Error("Expected some values to be loaded")
	}
}

func getFilm() filmV2 {
	apiFilm := filmV2{
		ID:      123,
		Title:   "Film One",
		Slug:    "/film/52",
		Tagline: "Tag line",
		Runtime: 123,
		Trailers: []struct {
			URL  string `json:"url"`
			Type string `json:"type"`
		}{{
			URL:  "https://cdn/trailer.mp4",
			Type: "Full",
		}},
		Cast: []struct {
			Name      string `json:"name"`
			Character string `json:"character"`
		}{{
			Name:      "James Earl Jones",
			Character: "Darth Vadar",
		}},
		Crew: []struct {
			Name string `json:"name"`
			Job  string `json:"job"`
		}{{
			Name: "Peter Jackson",
			Job:  "Director",
		}},
		Subtitles: []string{"Japanese"},
		SubtitleTracks: []subtitleTrackV1{{
			Language: "it",
			Name:     "Italian",
			Type:     "caption",
			Path:     "/subtitles/film/49/bonus/1/it/caption-18.vtt",
		}, {
			Language: "it",
			Name:     "Italian",
			Type:     "caption",
			Path:     "/subtitles/film/49/bonus/1/it/caption-19.vtt",
		}},
		Recommendations: []string{"/film/1", "/film/2"},
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
		Studio: []struct {
			Name string `json:"name"`
		}{{
			Name: "Studio ABC",
		}, {
			Name: "Studio XYZ",
		}},
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
		SeoTitle:       "Film One Meta Title",
		SeoKeywords:    "Film One Meta Keywords",
		SeoDescription: "Film One Meta Description",
	}
	return apiFilm
}

func TestFilmApiToModel(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()
	apiFilm := getFilm()

	model := apiFilm.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "film-one", model.TitleSlug, "title slug")
	assert.Equal(t, "/film/52", model.Slug, "slug")
	assert.Equal(t, "https://cdn/trailer.mp4", model.Trailers[0].URL, "trailer")

	assert.Equal(t, "SHIFT72 , Film One Meta Title , VOD", model.Seo.Title, "seo.title")
	assert.Equal(t, "Film One Meta Description", model.Seo.Description, "seo.description")
	assert.Equal(t, "SHIFT72, VOD, Film One Meta Keywords", model.Seo.Keywords, "seo.keywords")
	assert.Equal(t, "", model.Seo.Image, "seo.image")
	assert.Equal(t, "https://cdn/trailer.mp4", model.Seo.VideoURL, "seo.videourl")

	assert.Equal(t, "Darth Vadar", model.Cast[0].Character, "cast.character")
	assert.Equal(t, "Peter Jackson", model.Crew[0].Name, "crew.name")

	assert.Equal(t, 1, len(model.Bonuses), "expect 1 bonus")
	assert.Equal(t, "/film/52/bonus/1", model.Bonuses[0].Slug, "bonus.slug")

	assert.Equal(t, 2, len(model.Recommendations), "expect 2 generic items")

	assert.Equal(t, 2, len(itemIndex["film"]), "expect the item index to include 2 films")

	assert.Equal(t, 1, len(model.Subtitles), "expect hard-coded subs to be 1")
	assert.Equal(t, 2, len(model.SubtitleTracks), "expect the subtitles to be 2")

	assert.Equal(t, nil, model.CustomFields["hello?"])

	assert.Equal(t, 2, len(model.GetSubtitles()), "expect merged list of subtitles")
	assert.Contains(t, model.GetSubtitles(), "Italian")
	assert.Contains(t, model.GetSubtitles(), "Japanese")

	assert.Equal(t, 2, len(model.Classifications))

	assert.Equal(t, 2, len(model.Studio))
	assert.Contains(t, model.Studio, "Studio XYZ")
	assert.Contains(t, model.Studio, "Studio ABC")
}

func TestFilmApiToModelWithoutClassifications(t *testing.T) {
	itemIndex := make(models.ItemIndex)
	serviceConfig := commonServiceConfig()
	apiFilm := getFilm()
	apiFilm.Classifications = nil

	model := apiFilm.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 0, len(model.Classifications))
}

func TestFilmApiToModelWithoutSeoImage(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiFilm := getFilm()
	imageURL := "image.jpeg"
	apiFilm.ImageUrls.Portrait = imageURL
	apiFilm.ImageUrls.Landscape = imageURL

	model := apiFilm.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, model.Seo.Image, imageURL, "should be equal")
}

func TestFilmApiToModelWithSeoImage(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiFilm := getFilm()
	imageURL := "seo_image.jpeg"
	apiFilm.ImageUrls.Seo = imageURL

	model := apiFilm.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, model.Seo.Image, imageURL, "should be equal")
}

func TestFilmCustomFields(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiFilm := filmV2{
		ID:    123,
		Title: "Film One",
		Slug:  "/film/52",
		CustomFields: map[string]interface{}{
			"facebook_url": "https://www.facebook.com/custompage",
			"some_key":     1,
			"another_key":  false,
		},
	}

	model := apiFilm.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 1, model.CustomFields["some_key"])
	assert.Equal(t, "https://www.facebook.com/custompage", model.CustomFields["facebook_url"])
	assert.Equal(t, false, model.CustomFields["another_key"])
	assert.Equal(t, nil, model.CustomFields["where is it"])
}

func TestFilmSubtitlesAsArray(t *testing.T) {
	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiFilm := filmV2{
		ID:        123,
		Title:     "Film One",
		Slug:      "/film/52",
		Subtitles: []string{"French", "Italian", "French"},
	}

	model := apiFilm.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, 2, len(model.GetSubtitles()))
	assert.Contains(t, model.GetSubtitles(), "French")
	assert.Contains(t, model.GetSubtitles(), "Italian")
}

func TestUniqueFilmTitles(t *testing.T) {

	collection := models.FilmCollection{
		models.Film{
			ID:        1,
			Slug:      "/film/1",
			TitleSlug: "the-big-lebowski",
		},
		models.Film{
			ID:        2,
			Slug:      "/film/2",
			TitleSlug: "the-big-lebowski",
		},
	}

	collection.MakeTitleSlugsUnique()

	assert.Equal(t, "the-big-lebowski", collection[0].TitleSlug)
	assert.Equal(t, "the-big-lebowski-2", collection[1].TitleSlug)
}

func TestUniqueFilmTitlesTheNewestGetsTheLargestIndex(t *testing.T) {

	collection := models.FilmCollection{
		models.Film{
			ID:        3,
			Slug:      "/film/3",
			TitleSlug: "the-big-lebowski",
		},
		models.Film{
			ID:        1,
			Slug:      "/film/1",
			TitleSlug: "the-big-lebowski",
		},
		models.Film{
			ID:        2,
			Slug:      "/film/2",
			TitleSlug: "the-big-lebowski",
		},
	}

	collection.MakeTitleSlugsUnique()

	assert.Equal(t, "the-big-lebowski-3", collection[0].TitleSlug)
	assert.Equal(t, "the-big-lebowski", collection[1].TitleSlug)
	assert.Equal(t, "the-big-lebowski-2", collection[2].TitleSlug)
}
