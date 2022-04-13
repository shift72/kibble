package api

import (
	"testing"

	"kibble/models"

	"github.com/stretchr/testify/assert"
)

func TestBonusContentSubtitlesModelSupport(t *testing.T) {
	apiBonus := bonusContentV2{
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
	}

	model := apiBonus.mapToModel2("/film/1", models.ImageSet{
		Portrait: "film-portrait",
	})

	assert.Equal(t, 2, len(model.SubtitleTracks), "expect the subtitles to be 1")
	assert.Equal(t, "it", model.SubtitleTracks[0].Language)
	assert.Equal(t, "Italian", model.SubtitleTracks[0].Name)
	assert.Equal(t, "caption", model.SubtitleTracks[0].Type)
	assert.Equal(t, "/subtitles/film/49/bonus/1/it/caption-18.vtt", model.SubtitleTracks[0].Path)

	assert.Equal(t, 1, len(model.GetSubtitles()), "expect the subtitles to be 1")
	assert.Contains(t, model.GetSubtitles(), "Italian", "expect the subtitles to be the right name")
}

func TestBonusContentImagesUseFilmImagesAsFallback(t *testing.T) {

	itemIndex := models.NewItemIndex()

	serviceConfig := commonServiceConfig()

	apiFilm := filmV2{
		ID:      123,
		Title:   "Film One",
		Slug:    "/film/52",
		Tagline: "Tag line",
		Runtime: 123,
		Bonuses: []bonusContentV2{{
			Number: 1,
			Title:  "Behind the scenes",
		}},
		ImageUrls: map[string]string{
			"portrait_image":       "film-portrait.jpeg",
			"landscape_image":      "film-landscape.jpeg",
			"header_image":         "film-header.jpeg",
			"carousel_image":       "film-carousel.jpeg",
			"bg_image":             "film-background.jpeg",
			"classification_image": "film-classification.jpeg",
			"seo_image":            "film-seo.jpeg",
		},
	}

	model := apiFilm.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "film-portrait.jpeg", model.Bonuses[0].Images.Portrait)
	assert.Equal(t, "film-landscape.jpeg", model.Bonuses[0].Images.Landscape)
	assert.Equal(t, "film-header.jpeg", model.Bonuses[0].Images.Header)
	assert.Equal(t, "film-carousel.jpeg", model.Bonuses[0].Images.Carousel)
	assert.Equal(t, "film-background.jpeg", model.Bonuses[0].Images.Background)
	assert.Equal(t, "film-classification.jpeg", model.Bonuses[0].Images.Classification)
}

func TestBonusContentImagesUseSeasonImagesAsFallback(t *testing.T) {

	itemIndex := models.NewItemIndex()

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
		}},
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
	}

	model := apiSeason.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "season-portrait.jpeg", model.Bonuses[0].Images.Portrait)
	assert.Equal(t, "season-landscape.jpeg", model.Bonuses[0].Images.Landscape)
	assert.Equal(t, "season-header.jpeg", model.Bonuses[0].Images.Header)
	assert.Equal(t, "season-carousel.jpeg", model.Bonuses[0].Images.Carousel)
	assert.Equal(t, "season-background.jpeg", model.Bonuses[0].Images.Background)
	assert.Equal(t, "season-classification.jpeg", model.Bonuses[0].Images.Classification)
}

func TestBonusContentCustomFields(t *testing.T) {

	apiBonus := bonusContentV2{
		CustomFields: map[string]interface{}{
			"facebook_url": "https://www.facebook.com/custompage",
			"some_key":     1,
			"another_key":  false,
		},
	}

	model := apiBonus.mapToModel2("/film/1", models.ImageSet{})

	assert.Equal(t, 1, model.CustomFields["some_key"])
	assert.Equal(t, "https://www.facebook.com/custompage", model.CustomFields["facebook_url"])
	assert.Equal(t, false, model.CustomFields["another_key"])
	assert.Equal(t, nil, model.CustomFields["where is it"])

}
