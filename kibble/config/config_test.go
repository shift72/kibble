package config

import (
	"testing"
	"kibble/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadLanguagesConfigWithStrings(t *testing.T) {
	cfg := models.Config{}
	file := []byte(`{
		"languages": {
			"en": "en_AU",
			"it": "it_IT"
		}
	}`)

	LoadLanguagesConfig(&cfg, file)

	assert.Equal(t, cfg.Languages["en"].Code, "en_AU")
	assert.Equal(t, cfg.Languages["en"].Name, "")
	assert.Equal(t, cfg.Languages["it"].Code, "it_IT")
	assert.Equal(t, cfg.Languages["it"].Name, "")
}

func TestLoadLanguagesConfigWithObjects(t *testing.T) {
	cfg := models.Config{}
	file := []byte(`{
		"languages": {
			"en": {"code": "en_AU", "name": "English" },
			"it": {"code": "it_IT", "name": "Italian" }
		}
	}`)

	LoadLanguagesConfig(&cfg, file)

	assert.Equal(t, cfg.Languages["en"].Code, "en_AU")
	assert.Equal(t, cfg.Languages["en"].Name, "English")
	assert.Equal(t, cfg.Languages["it"].Code, "it_IT")
	assert.Equal(t, cfg.Languages["it"].Name, "Italian")
}

func TestLoadConfig(t *testing.T) {
	sitePath = "../sample_site/kibble.json"
	cfg := LoadConfig(false, "xxx", false)
	assert.Equal(t, true, cfg.RemoteTranslations)
}