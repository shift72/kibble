package api

import (
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
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

func TestFilmApiToModel(t *testing.T) {

	itemIndex := make(models.ItemIndex)

	serviceConfig := commonServiceConfig()

	apiFilm := filmV2{
		ID:      123,
		Title:   "Film One",
		Slug:    "/film/52",
		Tagline: "Tag line",
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
		Subtitles: []subtitleTrackV1{{
			Language: "it",
			Name:     "Italian",
			Type:     "caption",
			Path:     "/subtitles/film/49/bonus/1/it/caption-18.vtt",
		}},
		Recommendations: []string{"/film/1", "/film/2"},
		Bonuses: []filmBonusV2{{
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

	model := apiFilm.mapToModel(serviceConfig, itemIndex)

	assert.Equal(t, "film-one", model.TitleSlug, "title slug")
	assert.Equal(t, "/film/52", model.Slug, "slug")
	assert.Equal(t, "https://cdn/trailer.mp4", model.Trailers[0].URL, "trailer")

	assert.Equal(t, "SHIFT72 , Film One,  VOD", model.Seo.Title, "seo.title")
	assert.Equal(t, "Tag line", model.Seo.Description, "seo.description")
	assert.Equal(t, "SHIFT72, VOD", model.Seo.Keywords, "seo.keywords")
	assert.Equal(t, "", model.Seo.Image, "seo.image")
	assert.Equal(t, "https://cdn/trailer.mp4", model.Seo.VideoURL, "seo.videourl")

	assert.Equal(t, "Darth Vadar", model.Cast[0].Character, "cast.character")
	assert.Equal(t, "Peter Jackson", model.Crew[0].Name, "crew.name")

	assert.Equal(t, 1, len(model.Bonuses), "expect 1 bonus")
	assert.Equal(t, "/film/52/bonus/1", model.Bonuses[0].Slug, "bonus.slug")

	assert.Equal(t, 2, len(model.Recommendations), "expect 2 generic items")

	assert.Equal(t, 2, len(itemIndex["film"]), "expect the item index to include 2 films")

	assert.Equal(t, 1, len(model.Subtitles), "expect the subtitles to be 1")
}
