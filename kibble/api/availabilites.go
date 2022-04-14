package api

import (
	"context"
	"encoding/json"
	"fmt"
	"kibble/models"
	"kibble/utils"
	"sort"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

const batchSize = 300

// LoadAllAvailabilities will load all availabilities
func LoadAllAvailabilities(ctx context.Context, cfg *models.Config, site *models.Site, itemIndex *models.ItemIndex) error {
	sw := utils.NewStopwatchfWithLevel(" LoadAllAvailabilities ")
	slugs := make([]string, 0)

	for k := range itemIndex.Items["film"] {
		slugs = append(slugs, k)
	}

	for k := range itemIndex.Items["tv-season"] {
		slugs = append(slugs, k)
	}

	//NB: episode availabilities will be returned with the season

	sort.Strings(slugs)

	var total = 0
	g, _ := errgroup.WithContext(ctx)
	// create a channel to receive the results/no of items processed.
	res := make(chan int)
	for len(slugs) > 0 {
		var s []string

		if len(slugs) > batchSize {
			s = slugs[:batchSize]
			slugs = slugs[batchSize:]
		} else {
			s = slugs[:]
			slugs = nil
		}
		g.Go(func() error {
			return loadAvailabilities(cfg, site, s, itemIndex, res)
		})
	}
	go func() {
		_ = g.Wait()
		close(res)
	}()
	for count := range res {
		total += count
	}
	// Check whether any of the goroutines failed.
	if err := g.Wait(); err != nil {
		return err
	}
	log.Infof("availabilities: loaded %d", total)
	sw.Completed()
	return nil
}

func loadAvailabilities(cfg *models.Config, site *models.Site, slugs []string, itemIndex *models.ItemIndex, res chan int) error {
	ids := strings.Join(slugs, ",")
	path := fmt.Sprintf("%s/services/content/v1/availabilities?items=%s", cfg.SiteURL, ids)

	body, err := Get(cfg, path)
	if err != nil {
		log.Infof("availabilites failed to load %s", err)
		return err
	}

	var data availabilities
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("price.error: %s", err)
		log.Debug("invalid data %s", string(body))
		return err
	}

	c, err := processAvailabilities(data, site, itemIndex)
	if err != nil {
		log.Infof("processing availabilites failed %s", err)
		return err
	}
	res <- c

	return nil
}

func processAvailabilities(data availabilities, site *models.Site, itemIndex *models.ItemIndex) (int, error) {

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
		case "tvepisode":
			if f, ok := site.TVEpisodes.FindTVEpisodeBySlug(p.Slug); ok {
				count++
				f.Available = available
				// episodes do not appear in the itemIndex
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
