package models

import "testing"
import "github.com/stretchr/testify/assert"

func TestFilmSEO(t *testing.T) {

	serviceConfig := ServiceConfig{
		"image_root_path":       "https://s3-bla-bla",
		"portrait_poster_path":  "/posters-and-backdrops/282x422",
		"landscape_poster_path": "/posters-and-backdrops/380x210",
		"seo_title_prefix":      "SHIFT72 ",
		"seo_title_suffix":      " VOD",
		"seo_site_keywords":     "SHIFT72, VOD",
		"seo_site_name":         "Film On Demand",
	}

	film := &Film{
		Title:    "Brave",
		Keywords: "key, words",
		// description
	}

	seo := film.GetSEO(serviceConfig)

	assert.Equal(t, "SHIFT72 Brave VOD", seo.Title, "SEO")
	assert.Equal(t, "SHIFT72, VOD, key, words", seo.Keywords, "expect the site keywords and the film keywords to be joined")
}

func TestFilmSEODefaults(t *testing.T) {

	serviceConfig := ServiceConfig{
		"image_root_path":       "https://s3-bla-bla",
		"portrait_poster_path":  "/posters-and-backdrops/282x422",
		"landscape_poster_path": "/posters-and-backdrops/380x210",
		"seo_title_prefix":      "SHIFT72 ",
		"seo_title_suffix":      " VOD",
		"seo_site_keywords":     "SHIFT72, VOD",
		"seo_site_name":         "Film On Demand",
	}

	film := &Film{
		Title:    "Brave",
		Keywords: "",
		// description
	}

	seo := film.GetSEO(serviceConfig)

	assert.Equal(t, "SHIFT72 Brave VOD", seo.Title, "SEO")
	assert.Equal(t, "SHIFT72, VOD", seo.Keywords, "expect the site keywords")
}
