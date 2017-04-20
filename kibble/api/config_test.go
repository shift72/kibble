package api

import (
	"fmt"
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

func TestLoadConfig(t *testing.T) {

	if testing.Short() {
		return
	}

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	serviceConfig, err := LoadConfig(cfg)
	if err != nil {
		t.Error(err)
	}

	if len(serviceConfig) == 0 {
		t.Error("Expected some values to be loaded")
	}

	fmt.Printf("loaded service config: %d\n", len(serviceConfig))
}

func TestLoadFeatureToggles(t *testing.T) {

	if testing.Short() {
		return
	}

	cfg := &models.Config{
		SiteURL: "https://staging-store.shift72.com",
	}

	toggles, err := LoadFeatureToggles(cfg)
	if err != nil {
		t.Error(err)
	}

	if len(toggles) == 0 {
		t.Error("Expected some values to be loaded")
	}

	fmt.Printf("loaded feature toggles: %d\n", len(toggles))
}
