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

package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kibble/api"
	"kibble/models"
)

type LanguageRenderer interface {
	PreprocessLanguageFiles(sourcePath string) error
	OverrideLanguages() error
	ObtainTranslations() (api.TranslationsV1, error)
	FormatContextLanguage(translationsAPIEnabled bool, languageObjKey string, code string) *models.Language
	CreateLanguage(languageObjKey string, code string) *models.Language
}

type languageRenderer struct {
	cfg       *models.Config
	site      *models.Site
	isEnabled bool
}

func NewLanguageRenderer(cfg *models.Config, site *models.Site) LanguageRenderer {
	return &languageRenderer{
		cfg:       cfg,
		site:      site,
		isEnabled: site.Toggles["translations_api"],
	}
}

//Setup and Create language files based on API or local language files
func (l *languageRenderer) PreprocessLanguageFiles(sourcePath string) error {
	err := l.OverrideLanguages()
	if err != nil {
		return err
	}

	apiTranslations, err := l.ObtainTranslations()
	if err != nil {
		log.Errorf("Failed to get translations: %s", err)
		return err
	}

	if apiTranslations != nil {
		for _, languageObj := range l.cfg.Languages {
			code := languageObj.Code

			filename := formatLanguageFilename(code)

			file, err := json.Marshal(apiTranslations[code])
			if err != nil {
				log.Errorf("Failed to marshal translations json %s: %s", code, err)
				return err
			}

			err = writeFile(filepath.Join(sourcePath, filename), file)
			if err != nil {
				log.Errorf("Failed to write translations files: %s", err)
				return err
			}
		}
	}

	return nil
}

func (l *languageRenderer) OverrideLanguages() error {
	if !(l.isEnabled) {
		return nil
	}

	l.site.Languages = make([]models.Language, 0)

	l.cfg.DefaultLanguage = ""
	l.cfg.Languages = make(map[string]models.LanguageConfig)

	languages, err := api.LoadAllLanguages(l.cfg)
	if err != nil {
		log.Errorf("Failed to get languages: %s", err)
		return err
	}

	l.cfg.DefaultLanguage = formatPathLocale(languages.DefaultLanguage.Code)

	for _, lang := range languages.SupportedLanguages {
		langCode := formatPathLocale(lang.Code)
		langLabel := formatPathLocale(lang.Label)

		isDefault := langCode == l.cfg.DefaultLanguage

		l.site.Languages = append(l.site.Languages, models.Language{
			IsDefault: isDefault,
			Code:      defaultLanguageOverride(isDefault, langCode),
			// This is for selector display, will probably change after is(a)+c adds "Display Names"
			Name: lang.Label,
		})

		l.cfg.Languages[langCode] = models.LanguageConfig{
			Code: langCode,
			Name: langLabel,
		}
	}

	return nil
}

func (l *languageRenderer) ObtainTranslations() (api.TranslationsV1, error) {
	if !(l.isEnabled) {
		return nil, nil
	}

	translations, err := api.LoadAllTranslations(l.cfg)
	if err != nil {
		return nil, err
	}

	for key := range translations {
		translations[formatPathLocale(key)] = translations[key]
	}

	return translations, nil
}

func (l *languageRenderer) FormatContextLanguage(translationsAPIEnabled bool, languageObjKey string, code string) *models.Language {
	if translationsAPIEnabled {
		return &models.Language{
			Code:               formatPathLocale(languageObjKey),
			Locale:             formatPathLocale(code),
			IsDefault:          (formatPathLocale(languageObjKey) == l.cfg.DefaultLanguage),
			DefinitionFilePath: formatLanguageFilename(code),
		}
	}
	return l.CreateLanguage(languageObjKey, code)
}

func (l *languageRenderer) CreateLanguage(languageObjKey string, code string) *models.Language {
	return &models.Language{
		Code:               languageObjKey,
		Locale:             code,
		IsDefault:          (languageObjKey == l.cfg.DefaultLanguage),
		DefinitionFilePath: formatLanguageFilename(code),
	}
}

func formatPathLocale(code string) string {
	dashedCode := strings.ReplaceAll(code, "_", "-")
	return strings.ToLower(dashedCode)
}

func writeFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Errorf("%s", err)
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		log.Errorf("%s", err)
		file.Close()
		return err
	}

	return file.Close()
}

func formatLanguageFilename(code string) string {
	return fmt.Sprintf("%s.all.json", code)
}

func defaultLanguageOverride(isDefault bool, langCode string) string {
	if isDefault {
		return ""
	}
	return langCode
}
