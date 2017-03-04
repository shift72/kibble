package api

import (
	"fmt"
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

func TestLoadFilms(t *testing.T) {

	// if testing.Short() {
	// 	return
	// }

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	summary, err := LoadFilmSummary(cfg)
	if err != nil {
		t.Error(err)
	}

	if len(summary) == 0 {
		t.Error("Expected some values to be loaded")
	}

	fmt.Printf("loaded film summaries: %d\n", len(summary))
}

func TestLoadAllFilms(t *testing.T) {

	// if testing.Short() {
	// 	return
	// }

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	summary, err := LoadAllFilms(cfg)
	if err != nil {
		t.Error(err)
	}

	if len(summary) == 0 {
		t.Error("Expected some values to be loaded")
	}

	fmt.Printf("loaded film summaries: %d\n", len(summary))
}
