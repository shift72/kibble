package api

import (
	"encoding/json"
	"fmt"
	"kibble/models"
	"sort"
	"strings"
	"time"
)

// LoadAllAvailabilities will load all availabilities
func LoadAllAvailabilities(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	slugs := make([]string, 0)

	for k := range itemIndex["film"] {
		slugs = append(slugs, k)
	}

	for k := range itemIndex["tv-season"] {
		slugs = append(slugs, k)
	}

	sort.Strings(slugs)

	const batchSize = 300
	var total = 0
	var s []string
	for len(slugs) > 0 {

		if len(slugs) > batchSize {
			s = slugs[:batchSize]
			slugs = slugs[batchSize:]
		} else {
			s = slugs[:]
			slugs = nil
		}

		count, err := loadAvailabilities(cfg, site, s, itemIndex)
		if err != nil {
			return err
		}
		total += count
	}

	log.Infof("availabilities: loaded %d", total)

	return nil
}

func loadAvailabilities(cfg *models.Config, site *models.Site, slugs []string, itemIndex models.ItemIndex) (int, error) {

	ids := strings.Join(slugs, ",")
	path := fmt.Sprintf("%s/services/content/v1/availabilities?items=%s", cfg.SiteURL, ids)

	body, err := Get(cfg, path)
	if err != nil {
		log.Infof("pricing failed to load %s", err)
		return 0, err
	}

	var data availabilities
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Error("price.error: %s", err)
		log.Debug("invalid data %s", string(body))
		return 0, err
	}

	return processAvailabilities(data, site, itemIndex)
}

func processAvailabilities(data availabilities, site *models.Site, itemIndex models.ItemIndex) (int, error) {

	count := 0
	for _, p := range data {
		found := itemIndex.Get(p.Slug)

		available := models.Period{
			From: p.From,
			To:   p.To,
		}

		switch found.ItemType {
		case "film":
			if f, ok := site.Films.FindFilmBySlug(p.Slug); ok {
				count++
				f.Available = available
				// replace the itemIndex
				itemIndex.Replace(p.Slug, f.GetGenericItem())
			}
		case "tvseason":
			if f, ok := site.TVSeasons.FindTVSeasonBySlug(p.Slug); ok {
				count++
				f.Available = available
				// replace the itemIndex
				itemIndex.Replace(p.Slug, f.GetGenericItem())
			}
		}
	}

	return count, nil
}

type availabilities []struct {
	Slug string     `json:"slug"`
	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
	// RentalDurationMinutes         int        `json:"rental_duration_minutes"`
	// RentalPlaybackDurationMinutes int        `json:"rental_playback_duration_minutes"`
}
