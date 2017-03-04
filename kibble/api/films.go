package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gosimple/slug"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// LoadFilmSummary - load the bios request
func LoadFilmSummary(cfg *models.Config) ([]models.FilmSummary, error) {

	summary := []models.FilmSummary{}

	path := fmt.Sprintf("%s/services/meta/v2/film/index", cfg.SiteURL)

	data, err := Get(path)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(data))

	err = json.Unmarshal([]byte(data), &summary)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// LoadAllFilms -
func LoadAllFilms(cfg *models.Config) ([]models.Film, error) {

	summary, err := LoadFilmSummary(cfg)
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(summary))
	for i := 0; i < len(summary); i++ {
		ids[i] = summary[i].ID
	}
	return LoadFilmDetails(cfg, ids)
}

// LoadFilmDetails - load all films
func LoadFilmDetails(cfg *models.Config, filmIds []int) ([]models.Film, error) {

	sort.Ints(filmIds)
	ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(filmIds)), ","), "[]")

	path := fmt.Sprintf("%s/services/meta/v2/film/%s/show_multiple", cfg.SiteURL, ids)

	data, err := Get(path)
	if err != nil {
		return nil, err
	}

	details := []models.Film{}

	err = json.Unmarshal([]byte(data), &details)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(details); i++ {
		// fix subtitles
		if details[i].SubtitlesRaw == nil {
			// do nothing
		} else if strings.HasPrefix(string(details[i].SubtitlesRaw), "\"") {
			details[i].Subtitles = []string{string(details[i].SubtitlesRaw)}
		} else if strings.HasPrefix(string(details[i].SubtitlesRaw), "[") {
			details[i].Subtitles = strings.Split(strings.Trim(string(details[i].SubtitlesRaw), "[]"), ",")
		}

		details[i].TitleSlug = slug.Make(details[i].Title)
	}

	return details, nil
}
