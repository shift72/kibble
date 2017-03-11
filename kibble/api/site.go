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

	toggles, err := LoadFeatureToggles(cfg)
	if err != nil {
		return nil, err
	}

	bios, err := LoadBios(cfg, itemIndex)
	if err != nil {
		return nil, err
	}

	site := &models.Site{
		Config:     config,
		Toggles:    toggles,
		Navigation: bios.Navigation,
		Pages:      bios.Pages,
		Films:      make(models.FilmCollection, 0),
		Bundles:    make(models.BundleCollection, 0),
	}

	err = AppendAllFilms(cfg, site, itemIndex)
	if err != nil {
		return nil, err
	}

	err = LoadAllBundles(cfg, site, itemIndex)
	if err != nil {
		return nil, err
	}

	// while there are unresolved film slugs
	s := itemIndex.FindUnresolvedSlugs("film")
	for len(s) > 0 {
		AppendFilms(cfg, site, s, itemIndex)
		s = itemIndex.FindUnresolvedSlugs("film")
	}

	//TODO: while there are unresolved tv seasons

	fmt.Printf("service config:\t%d\n", len(config))
	fmt.Printf("toggles:\t%d\n", len(toggles))
	fmt.Printf("pages:\t\t%d\n", len(bios.Pages))
	fmt.Printf("films:\t\t%d\n", len(site.Films))
	fmt.Printf("bundles:\t%d\n", len(site.Bundles))

	stop := time.Now()
	fmt.Printf("-------------------------\nLoad completed: %s\n-------------------------\n", stop.Sub(start))

	// itemIndex.Print()
	itemIndex.PrintStats()

	site.LinkItems(itemIndex)

	return site, nil
}
