package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/nicksnyder/go-i18n/i18n"
)

// LoadConfig - loaded the config
func LoadConfig() *models.Config {

	file, err := ioutil.ReadFile("./site.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	var cfg models.Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		fmt.Printf("Config file parsing error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("url: %s\n", cfg.SiteURL)

	loadLanguages(&cfg)

	fmt.Printf("languages: %d\n", len(cfg.Languages))

	return &cfg
}

func loadLanguages(cfg *models.Config) {
	i18n.MustLoadTranslationFile(fmt.Sprintf("%s.all.json", cfg.Languages[cfg.DefaultLanguage]))

	for _, locale := range cfg.Languages {
		i18n.LoadTranslationFile(fmt.Sprintf("%s.all.json", locale))
	}
}
