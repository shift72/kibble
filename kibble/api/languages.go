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

func LoadAllLanguages(cfg *models.Config) (*LanguagesV1, error) {

	path := fmt.Sprintf("%s/services/users/v1/languages", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return nil, err
	}

	var languages LanguagesV1

	err = json.Unmarshal([]byte(data), &languages)
	if err != nil {
		return nil, err
	}

	return &languages, nil
}

type LanguagesV1 struct {
	DefaultLanguage    map[string]string   `json:"default_language"`
	SupportedLanguages []map[string]string `json:"supported_languages"`
}
