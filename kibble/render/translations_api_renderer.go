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
	"kibble/models"
	"os"
	"path/filepath"
)

// Ensures language files exist for languges sepcified via Translations API
func ensureLanguageFiles(site *models.Site, sourcePath string) error {
	// Create new language files (e.g. en_AU.all.json) from the data in site.Languages that has been sourced from the API.
	for _, language := range site.Languages {
		code := language.Code

		if code == "" {
			code = site.DefaultLanguage
		}

		filename := fmt.Sprintf("%s.all.json", code)
		file, err := json.Marshal(site.Translations[code])
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

	return nil
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
