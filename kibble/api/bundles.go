package api

import (
	"encoding/json"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// LoadAllBundles - load all bundles
func LoadAllBundles(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	path := fmt.Sprintf("%s/services/meta/v1/bundles", cfg.SiteURL)
	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	details := []models.Bundle{}
	err = json.Unmarshal([]byte(data), &details)
	if err != nil {
		return err
	}

	for i := 0; i < len(details); i++ {

		details[i].Slug = fmt.Sprintf("/bundle/%d", details[i].ID)
		details[i].TitleSlug = slug.Make(details[i].Title)

		site.Bundles = append(site.Bundles, details[i])
		itemIndex.Set(details[i].Slug, details[i].GetGenericItem())
	}

	return nil
}
