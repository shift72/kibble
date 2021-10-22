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

func LoadAllTranslations(cfg *models.Config, site *models.Site) error {
	if !site.Toggles["translations_api"] {
		return nil
	}

	path := fmt.Sprintf("%s/services/users/v1/translations", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	var translations TranslationsV1

	err = json.Unmarshal([]byte(data), &translations)
	if err != nil {
		return err
	}

	for code, wholeLanguage := range translations {
		site.Translations[formatPathLocale(code)] = wholeLanguage
	}

	return nil
}

type TranslationsV1 map[string]map[string]map[string]string
