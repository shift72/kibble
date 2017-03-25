package api

import (
	"encoding/json"
	"fmt"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

// LoadConfig - load all and merge
func LoadConfig(cfg *models.Config) (models.ServiceConfig, error) {

	var loaded models.ServiceConfig
	config := make(models.ServiceConfig)

	paths := []string{
		fmt.Sprintf("%s/services/users/v1/configuration", cfg.SiteURL),
		fmt.Sprintf("%s/services/shopping/configuration", cfg.SiteURL),
		fmt.Sprintf("%s/services/content/configuration", cfg.SiteURL),
		fmt.Sprintf("%s/services/meta/configuration", cfg.SiteURL),
	}

	for _, p := range paths {

		data, err := Get(cfg, p)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(data), &loaded)
		if err != nil {
			return nil, err
		}

		for k, v := range loaded {
			config[k] = v
		}
	}
	return config, nil
}

// LoadFeatureToggles - load all and merge
func LoadFeatureToggles(cfg *models.Config) (models.FeatureToggles, error) {

	var loaded models.FeatureToggles
	config := make(models.FeatureToggles)

	paths := []string{
		fmt.Sprintf("%s/services/shopping/feature_toggles", cfg.SiteURL),
		fmt.Sprintf("%s/services/meta/feature_toggles", cfg.SiteURL),
		fmt.Sprintf("%s/services/users/v1/feature_toggles", cfg.SiteURL),
	}

	for _, p := range paths {

		data, err := Get(cfg, p)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(data), &loaded)
		if err != nil {
			return nil, err
		}

		for k, v := range loaded {
			config[k] = v
		}
	}
	return config, nil
}
