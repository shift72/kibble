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
	"strings"
)

func LoadAllLanguages(cfg *models.Config, site *models.Site) error {
	if !site.Toggles["translations_api"] {
		return nil
	}

	path := fmt.Sprintf("%s/services/users/v1/languages", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return err
	}

	var languages languagesV1

	err = json.Unmarshal([]byte(data), &languages)
	if err != nil {
		return err
	}

	site.DefaultLanguage = formatPathLocale(languages.DefaultLanguage.Code)
	site.Languages = languages.mapToModel()

	return nil
}

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

func formatPathLocale(code string) string {
	dashedCode := strings.ReplaceAll(code, "_", "-")
	return strings.ToLower(dashedCode)
}
