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
	"encoding/json"
	"fmt"

	"kibble/models"
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
