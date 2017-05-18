package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	goversion "github.com/hashicorp/go-version"
	"github.com/indiereign/shift72-kibble/kibble/models"
	version "github.com/indiereign/shift72-kibble/kibble/version"
	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	privatePath       = path.Join(".kibble", "private.json")
	sitePath          = path.Join("site.json")
	currentVersion, _ = goversion.NewVersion(version.Version)
)

// LoadConfig loads the configuration from disk. If runAsAdmin it will attempt
// to load the private config
func LoadConfig(runAsAdmin, disableCache bool) *models.Config {

	file, err := ioutil.ReadFile(sitePath)
	if err != nil {
		log.Errorf("file error: %v\n", err)
		os.Exit(1)
	}

	var cfg models.Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Errorf("config file parsing error: %v", err)
		os.Exit(1)
	}

	cfg.RunAsAdmin = runAsAdmin
	cfg.DisableCache = disableCache

	log.Debugf("url: %s", cfg.SiteURL)

	loadLanguages(&cfg)

	LoadPrivateConfig(&cfg)

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

	if cfg.BuiltWithVersion == "" {
		return
	}

	bwv, err := goversion.NewVersion(cfg.BuiltWithVersion)
	if err != nil {
		log.Warning("invalid version, assuming version 0.0.0")
		bwv, _ = goversion.NewVersion("0.0.0")
	}

	if bwv.GreaterThan(currentVersion) {
		log.Warning("this template was built with a newer version of kibble, some templates possibly will not work")
	}
}

// UpdateBuiltWithVersion updates the build with version with the current version and saves the config
func UpdateBuiltWithVersion(cfg *models.Config) {
	if cfg.BuiltWithVersion != version.Version {
		cfg.BuiltWithVersion = version.Version
		SaveConfig(cfg)
	}
}

// LoadPrivateConfig is responsible for loading the private configuration if it exists
func LoadPrivateConfig(cfg *models.Config) {

	if !cfg.RunAsAdmin {
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

func loadLanguages(cfg *models.Config) {
	i18n.MustLoadTranslationFile(fmt.Sprintf("%s.all.json", cfg.Languages[cfg.DefaultLanguage]))

	for _, locale := range cfg.Languages {
		i18n.LoadTranslationFile(fmt.Sprintf("%s.all.json", locale))
	}
}
