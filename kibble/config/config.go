package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	goversion "github.com/hashicorp/go-version"
	"github.com/indiereign/shift72-kibble/kibble/models"
	version "github.com/indiereign/shift72-kibble/kibble/version"
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

	log.Debugf("url: %s", cfg.SiteURL)

	LoadPrivateConfig(&cfg, apiKey)

	return &cfg
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
func CheckVersion(cfg *models.Config) {

	if cfg.BuilderVersion == "" {
		return
	}

	bwv, err := goversion.NewVersion(cfg.BuilderVersion)
	if err != nil {
		log.Warning("invalid version, assuming version 0.0.0")
		bwv, _ = goversion.NewVersion("0.0.0")
	}

	if bwv.GreaterThan(currentVersion) {
		log.Warning("this template was built with a newer version of kibble, some templates possibly will not work")
	}
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
