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
	"sort"
	"strings"
)

//Loads all languages from the API if the site_translations_api feature toggle is enabled.
func LoadAllLanguages(cfg *models.Config, site *models.Site) error {
	if site.Toggles["site_translations_api"] {
		return loadAllLanguagesFromApi(cfg, site)
	} else {
		return loadAllLanguagesFromConfig(cfg, site)
	}
}

func loadAllLanguagesFromApi(cfg *models.Config, site *models.Site) error {
	path := fmt.Sprintf("%s/services/users/v1/languages", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return fmt.Errorf("languages from API failed to load %s", err)
	}

	var languages languagesV1

	err = json.Unmarshal([]byte(data), &languages)
	if err != nil {
		return err
	}

	err = languages.validate()
	if err != nil {
		return fmt.Errorf("Language validation failed: %s", err)
	}

	// Force default language to be lowercase with '-' instead of '_'.
	site.DefaultLanguage = formatPathLocale(languages.DefaultLanguage.Code)
	site.Languages = languages.mapToModel()

	log.Infof("Site Languages Loaded from API")
	var languagesList []string
	for _, langs := range site.Languages {
		if langs.IsDefault {
			languagesList = append(languagesList, site.DefaultLanguage)
		} else {
			languagesList = append(languagesList, langs.Code)
		}
	}
	log.Infof("Default Language: %s ", site.DefaultLanguage)
	log.Infof("Supported Languages: %s", strings.Join(languagesList, ", "))
	return nil
}

// Returns an error if:
// - default language is not set.
// - default language is missing display name.
// - no site languages are set.
// - any site language is missing a display name.
// - default language is not included in supported languages.
func (l languagesV1) validate() error {
	isZero := func(val string) bool {
		return val == ""
	}

	if isZero(l.DefaultLanguage.Code) && isZero(l.DefaultLanguage.Label) && isZero(l.DefaultLanguage.Name) {
		return fmt.Errorf("DefaultLanguage not set")
	}

	if isZero(l.DefaultLanguage.Name) {
		return fmt.Errorf("DefaultLanguage %s is missing display name", l.DefaultLanguage.Code)
	}

	if len(l.SupportedLanguages) == 0 {
		return fmt.Errorf("No SiteLanguages set")
	}

	defaultInLanguages := false
	for _, lang := range l.SupportedLanguages {
		if lang.Code == l.DefaultLanguage.Code {
			defaultInLanguages = true
		}

		if isZero(lang.Name) {
			return fmt.Errorf("Language %s is missing display name", lang.Code)
			// no need to check Code or Label as they are required by database schema.
		}
	}
	if !defaultInLanguages {
		return fmt.Errorf("DefaultLanguage not in SiteLanguages")
	}

	return nil
}

// Maps API response to array of Language models.
// Forces all language codes to be lowercase with '-' instead of '_' for use in browser path.
// Forces default language to have a blank code for use in browser path.
func (l languagesV1) mapToModel() []models.Language {
	languages := make([]models.Language, 0)

	for _, lang := range l.SupportedLanguages {
		code := formatPathLocale(lang.Code)
		isDefault := code == formatPathLocale(l.DefaultLanguage.Code)
		if isDefault {
			code = ""
		}
		languages = append(languages, models.Language{
			Code:      code,
			Name:      lang.Name,
			Label:     lang.Label,
			IsDefault: isDefault,
		})
	}
	return languages
}

type languagesV1 struct {
	DefaultLanguage    languageV1   `json:"default_language"`
	SupportedLanguages []languageV1 `json:"supported_languages"`
}

type languageV1 struct {
	Code  string `json:"code"`
	Name  string `json:"display_name"`
	Label string `json:"label"`
}

// Replace '_' with '-' and lowercase
func formatPathLocale(code string) string {
	dashedCode := strings.ReplaceAll(code, "_", "-")
	dashedCode = strings.TrimSpace(dashedCode)
	return strings.ToLower(dashedCode)
}

func loadAllLanguagesFromConfig(cfg *models.Config, site *models.Site) error {
	var keys []string
	for k := range cfg.Languages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		isDefault := k == cfg.DefaultLanguage
		code := k
		name := cfg.Languages[k].Name
		if isDefault {
			code = ""
		}

		site.Languages = append(site.Languages, models.Language{
			IsDefault: isDefault,
			Code:      code,
			Name:      name,
		})
	}
	return nil
}
