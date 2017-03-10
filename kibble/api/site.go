package api

import (
	"fmt"
	"time"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

// LoadSite - load the complete site
func LoadSite(cfg *models.Config) (*models.Site, error) {

	start := time.Now()

	itemIndex := make(models.ItemIndex)

	config, err := LoadConfig(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Printf("service config: %d\n", len(config))

	toggles, err := LoadFeatureToggles(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Printf("toggles: %d\n", len(toggles))

	bios, err := LoadBios(cfg, itemIndex)
	if err != nil {
		return nil, err
	}

	fmt.Printf("pages: %d\n", len(bios.Pages))

	films, err := LoadAllFilms(cfg, itemIndex)
	if err != nil {
		return nil, err
	}

	fmt.Printf("films: %d\n", len(films))

	stop := time.Now()
	fmt.Printf("--------------------\nLoad completed: %s\n--------------------\n", stop.Sub(start))

	site := &models.Site{
		Config:     config,
		Toggles:    toggles,
		Navigation: bios.Navigation,
		Pages:      bios.Pages,
		Films:      films,
	}

	site.IndexItems(itemIndex)

	// show the loaded items
	itemIndex.PrintStats()
	//TODO: request missing items

	site.LinkItems(itemIndex)

	return site, nil
}
