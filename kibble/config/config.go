package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/indiereign/shift72-kibble/kibble/models"
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

	fmt.Println("config loaded")

	return &cfg
}
