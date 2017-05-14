package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	privatePath = path.Join(".kibble", "private.json")
	sitePath    = path.Join("site.json")
)

// LoadConfig - loaded the config
func LoadConfig(runAsAdmin bool) *models.Config {

	file, err := ioutil.ReadFile(sitePath)
	if err != nil {
		log.Errorf("File error: %v\n", err)
		os.Exit(1)
	}

	var cfg models.Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Errorf("Config file parsing error: %v", err)
		os.Exit(1)
	}

	log.Debugf("url: %s", cfg.SiteURL)

	loadLanguages(&cfg)

	if runAsAdmin {
		LoadPrivateConfig(&cfg)
	}
	return &cfg
}

func SaveConfig(cfg *models.Config) {

	data, err := json.Marshal(cfg)
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

// LoadPrivateConfig - load any private configuratio
func LoadPrivateConfig(cfg *models.Config) {

	_, err := os.Stat(privatePath)
	if os.IsNotExist(err) {
		return
	}

	file, err := ioutil.ReadFile(privatePath)
	if err != nil {
		log.Errorf("File error: %v", err)
		os.Exit(1)
	}

	var private models.PrivateConfig
	err = json.Unmarshal(file, &private)
	if err != nil {
		log.Errorf("Config file parsing error: %v", err)
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
