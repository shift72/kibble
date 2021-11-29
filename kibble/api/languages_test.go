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
	"kibble/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLangaugeV1MapToModel(t *testing.T) {
	var languages languagesV1
	var body = "{\"default_language\":{\"code\":\"en_AU\",\"label\":\"English\",\"display_name\":\"EnglishDisplayName\"},\"supported_languages\":[{\"code\":\"de_DE\",\"label\":\"Deutsche\",\"display_name\":\"DeutscheDisplayName\"},{\"code\":\"en_AU\",\"label\":\"English\",\"display_name\":\"EnglishDisplayName\"}]}"
	var err = json.Unmarshal([]byte(body), &languages)
	if err != nil {
		t.Error(err)
	}
	var maptoModel = languages.mapToModel()
	assert.Equal(t, maptoModel, []models.Language{
		{Code: "de-de", Name: "DeutscheDisplayName", Label: "Deutsche", Locale: "", DefinitionFilePath: "", IsDefault: false},
		{Code: "", Name: "EnglishDisplayName", Label: "English", Locale: "", DefinitionFilePath: "", IsDefault: true}})
}
func TestLoadAllLanguagesFromConfig(t *testing.T) {

	cfg := &models.Config{
		SiteURL:         "https://staging-store.shift72.com",
		DefaultLanguage: "en",
		Languages: map[string]models.LanguageConfig{
			"en": {Code: "en-AU", Name: "English"},
			"de": {Code: "de-DE", Name: "Deutsche"},
			"fr": {Code: "fr-FR", Name: "French"},
		},
	}

	site := &models.Site{}
	loadAllLanguagesFromConfig(cfg, site)
	assert.Equal(t, site.Languages, []models.Language{
		{Code: "de", Name: "Deutsche", Label: "", Locale: "", DefinitionFilePath: "", IsDefault: false},
		{Code: "", Name: "English", Label: "", Locale: "", DefinitionFilePath: "", IsDefault: true},
		{Code: "fr", Name: "French", Label: "", Locale: "", DefinitionFilePath: "", IsDefault: false},
	})

}

func TestFormatPathLocale(t *testing.T) {
	var languageCode = "en_AU"
	var formatted = formatPathLocale(languageCode)
	assert.Equal(t, formatted, "en-au")

	languageCode = "en-au"
	formatted = formatPathLocale(languageCode)
	assert.Equal(t, formatted, "en-au")

	languageCode = "EN-AU"
	formatted = formatPathLocale(languageCode)
	assert.Equal(t, formatted, "en-au")

	languageCode = "ABC"
	formatted = formatPathLocale(languageCode)
	assert.Equal(t, formatted, "abc")

	languageCode = "123"
	formatted = formatPathLocale(languageCode)
	assert.Equal(t, formatted, "123")

	languageCode = "_"
	formatted = formatPathLocale(languageCode)
	assert.Equal(t, formatted, "-")

}
