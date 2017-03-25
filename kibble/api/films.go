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

	data, err := Get(cfg, path)
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

	fmt.Println("loading film count:", len(slugs))

	// set index to empty for the items requested
	for _, s := range slugs {
		itemIndex.Set(s, models.Empty)
	}

	path := fmt.Sprintf("%s/services/meta/v2/film/%s/show_multiple", cfg.SiteURL, ids)
	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	var details []json.RawMessage
	err = json.Unmarshal([]byte(data), &details)
	if err != nil {
		fmt.Println("film.error", err)
		fmt.Println(string(data))
		return err
	}

	for i := 0; i < len(details); i++ {

		var film models.Film
		err = json.Unmarshal(details[i], &film)

		if err == nil {

			// fix subtitles
			if film.SubtitlesRaw == nil {
				// do nothing
			} else if strings.HasPrefix(string(film.SubtitlesRaw), "\"") {
				film.Subtitles = []string{string(film.SubtitlesRaw)}
			} else if strings.HasPrefix(string(film.SubtitlesRaw), "[") {
				film.Subtitles = strings.Split(strings.Trim(string(film.SubtitlesRaw), "[]"), ",")
			}
			film.TitleSlug = slug.Make(film.Title)

			// add film
			site.Films = append(site.Films, film)
			itemIndex.Set(film.Slug, film.GetGenericItem())

			// add bonuses - supports linking to bonus entries (supported??)
			for bonusIndex := 0; bonusIndex < len(film.Bonuses); bonusIndex++ {
				itemIndex.Set(fmt.Sprintf("%s/bonus/%d", film.Slug, film.Bonuses[bonusIndex].Number), film.Bonuses[bonusIndex].GetGenericItem())
			}

			// add Recommendations
			for _, slug := range film.Recommendations {
				itemIndex.Set(slug, models.Unresolved)
			}

		} else {
			fmt.Println("err", err)
			fmt.Println("invalid data", string(details[i]))
		}

	}

	return nil
}
