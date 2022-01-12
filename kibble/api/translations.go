//    Copyright 2018 SHIFT72
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

// Loads all translations from the API if the site_translations_api feature toggle is enabled.
func LoadAllTranslations(cfg *models.Config, site *models.Site) error {
	if !site.Toggles["site_translations_api"] {
		// Do nothing - use translations aleady loaded in JSON files
		return nil
	}
	//Load from API
	path := fmt.Sprintf("%s/services/users/v1/site_translations", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return fmt.Errorf(" translations from API failed to load %s", err)
	}

	var translations TranslationsV1

	err = json.Unmarshal([]byte(data), &translations)
	if err != nil {
		return err
	}
	// Map translations to lowercase codes on the site
	for code, wholeLanguage := range translations {
		site.Translations[formatPathLocale(code)] = wholeLanguage
	}

	log.Infof("Translations Recieved and Parsed")
	return nil
}

// { "en-au": { "nav_signin": { "other": "Sign In" } } }
type TranslationsV1 map[string]map[string]struct {
	Zero  string `json:"zero,omitempty"`
	One   string `json:"one,omitempty"`
	Two   string `json:"two,omitempty"`
	Few   string `json:"few,omitempty"`
	Many  string `json:"many,omitempty"`
	Other string `json:"other,omitempty"`
}
