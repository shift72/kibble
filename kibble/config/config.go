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

package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	goversion "github.com/hashicorp/go-version"
	"kibble/models"
	version "kibble/version"
)

var (
	privatePath       = path.Join(".kibble", "private.json")
	sitePath          = path.Join("kibble.json")
	currentVersion, _ = goversion.NewVersion(version.Version)
)

// LoadConfig loads the configuration from disk. If runAsAdmin it will attempt
// to load the private config
func LoadConfig(runAsAdmin bool, apiKey string, disableCache bool) *models.Config {

	file, err := ioutil.ReadFile(sitePath)
	if err != nil {
		log.Errorf("file error: %v\n", err)
		os.Exit(1)
	}

	cfg := models.Config{
		RunAsAdmin:   runAsAdmin,
		DisableCache: disableCache,
	}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Errorf("config file parsing error: %v", err)
		os.Exit(1)
	}

	LoadLanguagesConfig(&cfg, file)

	log.Debugf("url: %s", cfg.SiteURL)

	LoadPrivateConfig(&cfg, apiKey)

	return &cfg
}

type languagesConfigObjects struct {
	Languages map[string]models.LanguageConfig `json:"languages"`
}
type languagesConfigStrings struct {
	Languages map[string]string `json:"languages"`
}

func LoadLanguagesConfig(cfg *models.Config, file []byte) {

	langCfgObj := languagesConfigObjects{}
	err := json.Unmarshal(file, &langCfgObj)
	if err != nil {

		langCfgStr := languagesConfigStrings{}
		err = json.Unmarshal(file, &langCfgStr)
		if err != nil {
			log.Errorf("config file languages parsing error: %v", err)
			os.Exit(1)
		}

		cfg.Languages = map[string]models.LanguageConfig{}

		for k, v := range langCfgStr.Languages {
			cfg.Languages[k] = models.LanguageConfig{Code: v}
		}

	} else {
		cfg.Languages = langCfgObj.Languages
	}
}

// SaveConfig writes the configuration to disk
func SaveConfig(cfg *models.Config) {

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Errorf("File error: %v\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(sitePath, data, 0777)
	if err != nil {
		log.Errorf("File error: %v", err)
		os.Exit(1)
	}
}

// CheckVersion is responsible for displaying an warning message if the version of kibble is newer than
// the version the template was built with
func CheckVersion(cfg *models.Config) error {

	if cfg.BuilderVersion == "" {
		log.Warning("No builderVersion specified.")
		return errors.New("No builderVersion specified")
	}

	bwv, err := goversion.NewVersion(cfg.BuilderVersion)
	if err != nil {
		log.Warning("invalid version, assuming version 0.0.0")
		bwv, _ = goversion.NewVersion("0.0.0")
	}

	if bwv.GreaterThan(currentVersion) {
		log.Warning("This template currently targets version %s, but you are running an older version %s.", cfg.BuilderVersion, currentVersion)
		log.Warning("Some features may not work, it is recommended that you update your version of kibble via npm.")
		return errors.New("Miss matched builderVersion")
	} else if bwv.LessThan(currentVersion) {
		log.Warning("This template currently targets version %s, but you are running an newer version %s.", cfg.BuilderVersion, currentVersion)
		log.Warning("It is recommended that you update the builderVersion in the kibble.json and test throughly before publishing.")
		return errors.New("Miss matched builderVersion")
	}

	return nil
}

// UpdateBuilderVersion updates the build with version with the current version and saves the config
func UpdateBuilderVersion(cfg *models.Config) {
	if cfg.BuilderVersion != version.Version {
		cfg.BuilderVersion = version.Version
		SaveConfig(cfg)
	}
}

// LoadPrivateConfig is responsible for loading the private configuration if it exists
func LoadPrivateConfig(cfg *models.Config, apiKey string) {

	if !cfg.RunAsAdmin {
		return
	}

	if apiKey != "" {
		cfg.SkipLogin = true
		cfg.Private = models.PrivateConfig{
			APIKey: apiKey,
		}
		return
	}

	_, err := os.Stat(privatePath)
	if os.IsNotExist(err) {
		return
	}

	file, err := ioutil.ReadFile(privatePath)
	if err != nil {
		log.Errorf("file error: %v", err)
		os.Exit(1)
	}

	var private models.PrivateConfig
	err = json.Unmarshal(file, &private)
	if err != nil {
		log.Errorf("config file parsing error: %v", err)
		os.Exit(1)
	}

	cfg.Private = private
}

// SavePrivateConfig - saves any private config
func SavePrivateConfig(cfg *models.Config) {

	data, err := json.Marshal(cfg.Private)
	if err != nil {
		log.Errorf("File error: %v", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(privatePath, data, 0777)
	if err != nil {
		log.Errorf("File error: %v", err)
		os.Exit(1)
	}
}
