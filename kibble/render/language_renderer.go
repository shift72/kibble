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
	ObtainTranslations() (api.TranslationsV1, error)
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
	apiTranslations, err := l.ObtainTranslations()
	if err != nil {
		log.Errorf("Failed to get translations: %s", err)
		return err
	}

	//Create translation filenames based on langague code
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

func (l *languageRenderer) ObtainTranslations() (api.TranslationsV1, error) {
	if !(l.isEnabled) {
		return nil, nil
	}

	//retrieve translations
	translations, err := api.LoadAllTranslations(l.cfg)
	if err != nil {
		return nil, err
	}

	for key := range translations {
		translations[formatPathLocale(key)] = translations[key]
	}

	return translations, nil
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
