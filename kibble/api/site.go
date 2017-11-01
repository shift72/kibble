package api

import (
	"sort"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/utils"
	logging "github.com/op/go-logging"
)

// LoadSite - load the complete site
func LoadSite(cfg *models.Config) (*models.Site, error) {

	initAPI := utils.NewStopwatchLevel("api", logging.NOTICE)

	itemIndex := make(models.ItemIndex)

	config, err := LoadConfig(cfg)
	if err != nil {
		return nil, err
	}

	toggles, err := LoadFeatureToggles(cfg)
	if err != nil {
		return nil, err
	}

	pages, navigation, err := LoadBios(cfg, config, itemIndex)
	if err != nil {
		return nil, err
	}

	site := &models.Site{
		SiteConfig:  cfg,
		Config:      config,
		Toggles:     toggles,
		Languages:   sortLanguages(cfg),
		Navigation:  navigation,
		Pages:       pages,
		Films:       make(models.FilmCollection, 0),
		Bundles:     make(models.BundleCollection, 0),
		Collections: make(models.CollectionCollection, 0),
		Taxonomies:  make(models.Taxonomies),
		TVSeasons:   make(models.TVSeasonCollection, 0),
	}

	err = LoadAllCollections(cfg, site, itemIndex)
	if err != nil {
		return nil, err
	}

	err = AppendAllFilms(cfg, site, itemIndex)
	if err != nil {
		return nil, err
	}

	err = LoadAllBundles(cfg, site, itemIndex)
	if err != nil {
		return nil, err
	}

	err = AppendAllTVSeasons(cfg, site, itemIndex)
	if err != nil {
		return nil, err
	}

	// while there are unresolved film slugs
	s := itemIndex.FindUnresolvedSlugs("film")
	for len(s) > 0 {
		AppendFilms(cfg, site, s, itemIndex)
		s = itemIndex.FindUnresolvedSlugs("film")
	}

	// while there are unresolved tv season slugs
	tvs := itemIndex.FindUnresolvedSlugs("tv-season")
	for len(tvs) > 0 {
		AppendTVSeasons(cfg, site, tvs, itemIndex)
		tvs = itemIndex.FindUnresolvedSlugs("tv-season")
	}

	initAPI.Completed()

	site.LinkItems(itemIndex)

	site.PopulateTaxonomyWithFilms("year", models.GetYear)
	site.PopulateTaxonomyWithFilms("genre", models.GetGenres)
	site.PopulateTaxonomyWithFilms("cast", models.GetCast)
	site.PopulateTaxonomyWithFilms("country", models.GetCountries)

	site.PopulateTaxonomyWithTVSeasons("year", models.GetTVSeasonYear)
	site.PopulateTaxonomyWithTVSeasons("genre", models.GetTVShowGenres)
	site.PopulateTaxonomyWithTVSeasons("cast", models.GetTVShowCast)
	site.PopulateTaxonomyWithTVSeasons("country", models.GetTVShowCountries)

	if log.IsEnabledFor(logging.DEBUG) {
		itemIndex.PrintStats()
	}

	return site, nil
}

func sortLanguages(cfg *models.Config) []models.Language {

	result := make([]models.Language, 0)

	var keys []string
	for k := range cfg.Languages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		result = append(result, models.Language{
			IsDefault: k == cfg.DefaultLanguage,
			Code:      k,
		})
	}
	return result
}
