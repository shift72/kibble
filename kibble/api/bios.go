package api

import (
	"encoding/json"
	"fmt"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

// LoadBios - load the bios request
func LoadBios(cfg *models.Config, itemIndex models.ItemIndex) (*models.Bios, error) {

	bios := &models.Bios{}

	path := fmt.Sprintf("%s/services/meta/v1/bios", cfg.SiteURL)

	data, err := Get(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &bios)
	if err != nil {
		return nil, err
	}

	// register with the item index
	for _, p := range bios.Pages {
		itemIndex.Add(fmt.Sprintf("/page/%d", p.ID), p.GetGenericItem())
	}

	return bios, nil
}
