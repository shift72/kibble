package api

import (
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

func TestLoadTVSeasons(t *testing.T) {

	// if testing.Short() {
	// 	return
	// }

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	itemIndex := make(models.ItemIndex)
	site := &models.Site{}
	slugs := []string{
		"/tv/4/season/2",
		"/tv/41/season/1",
		"/tv/9/season/1",
	}

	err := AppendTVSeasons(cfg, site, slugs, itemIndex)
	if err != nil {
		t.Error(err)
	}

	if len(itemIndex) == 0 {
		t.Error("Expected some values to be loaded")
	}
}
