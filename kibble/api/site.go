//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package api

import (
	"sort"

	"kibble/models"
	"kibble/utils"
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
		SiteConfig: cfg,
		Config:     config,
		Toggles:    toggles,
		Languages:  sortLanguages(cfg),
		Navigation: navigation,

		Bundles:     make(models.BundleCollection, 0),
		Collections: make(models.CollectionCollection, 0),
		Films:       make(models.FilmCollection, 0),
		Pages:       pages,
		Plans:       make(models.PlanCollection, 0),
		Taxonomies:  make(models.Taxonomies),
		TVShows:     make(models.TVShowCollection, 0),
		TVSeasons:   make(models.TVSeasonCollection, 0),
		TVEpisodes:  make(models.TVEpisodeCollection, 0),
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

	err = LoadAllPlans(cfg, site, itemIndex)
	if err != nil {
		return nil, err
	}

	err = AppendAllTVShows(cfg, site, itemIndex)
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
