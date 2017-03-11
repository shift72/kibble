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

	err = json.Unmarshal([]byte(data), &summary)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// AppendAllFilms -
func AppendAllFilms(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	summary, err := LoadFilmSummary(cfg)
	if err != nil {
		return err
	}

	for i := 0; i < len(summary); i++ {
		itemIndex.Set(fmt.Sprintf("/film/%d", summary[i].ID), models.Unresolved)
	}

	return nil
}

// AppendFilms - load a list of films
func AppendFilms(cfg *models.Config, site *models.Site, slugs []string, itemIndex models.ItemIndex) error {

	sort.Strings(slugs)

	// convert /film/1,film/2 -> 1,2
	ids := strings.Replace(
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(slugs)), ","), "[]"),
		"/film/", "", -1)

	// set index to empty for the items requested
	for _, s := range slugs {
		itemIndex.Set(s, models.Empty)
	}

	path := fmt.Sprintf("%s/services/meta/v2/film/%s/show_multiple", cfg.SiteURL, ids)
	data, err := Get(path)
	if err != nil {
		return err
	}

	details := []models.Film{}
	err = json.Unmarshal([]byte(data), &details)
	if err != nil {
		return err
	}

	for i := 0; i < len(details); i++ {

		// fmt.Println("film loaded", details[i].ID)

		// fix subtitles
		if details[i].SubtitlesRaw == nil {
			// do nothing
		} else if strings.HasPrefix(string(details[i].SubtitlesRaw), "\"") {
			details[i].Subtitles = []string{string(details[i].SubtitlesRaw)}
		} else if strings.HasPrefix(string(details[i].SubtitlesRaw), "[") {
			details[i].Subtitles = strings.Split(strings.Trim(string(details[i].SubtitlesRaw), "[]"), ",")
		}
		details[i].TitleSlug = slug.Make(details[i].Title)

		// add film
		site.Films = append(site.Films, details[i])
		itemIndex.Set(details[i].Slug, details[i].GetGenericItem())

		// add Recommendations
		for _, slug := range details[i].Recommendations {
			itemIndex.Set(slug, models.Unresolved)
		}

	}

	return nil
}
