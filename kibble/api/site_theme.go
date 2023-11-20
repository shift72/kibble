//    Copyright 2018 SHIFT72
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

func LoadSiteTheme(cfg *models.Config, site *models.Site) error {
	themeEnabled := site.Toggles["self_service_theme"]
	var siteTheme SiteThemeV1

	//If themes disabled
	if !themeEnabled {
		// Do nothing - use local site assets
		log.Infof("Themes disabled, using legacy branding api")
		return nil
	}

	//Load from API
	// for dev purpose https://core-store.shift72.local/services/users/v1/site_theme
	path := "https://core-store.shift72.local/services/users/v1/site_theme"
	// path := fmt.Sprintf("%s/services/users/v1/site_brand", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return fmt.Errorf(" Site theme from API failed to load %s", err)
	}

	err = json.Unmarshal([]byte(data), &siteTheme)
	if err != nil {
		return err
	}

	log.Infof("Themes enabled, getting theme data")
	theme := siteTheme.mapToModel(siteTheme)
	log.Infof("Theme loaded: %s", theme)
	log.Infof("primary: %s", theme.PrimaryRGB)
	log.Infof("secondary: %s", theme.SecondaryRGB)
	log.Infof("bg: %s", theme.BackgroundRGB)
	log.Infof("txt: %s", theme.TextRGB)
	log.Infof("logo size: %s", theme.LogoSize)
	site.SiteTheme = theme
	log.Infof("Theme loaded into site model:")
	log.Infof("primary: %s", site.SiteTheme.PrimaryRGB)
	log.Infof("secondary: %s", site.SiteTheme.SecondaryRGB)
	log.Infof("bg: %s", site.SiteTheme.BackgroundRGB)
	log.Infof("txt: %s", site.SiteTheme.TextRGB)
	log.Infof("logo size: %s", site.SiteTheme.LogoSize)

	return nil
}

func (t SiteThemeV1) mapToModel(siteTheme SiteThemeV1) models.SiteTheme {
	theme := models.SiteTheme{
		PrimaryRGB:    t.PrimaryRGB,
		SecondaryRGB:  t.SecondaryRGB,
		BackgroundRGB: t.BackgroundRGB,
		TextRGB:       t.TextRGB,
		LogoSize:      t.LogoSize,
	}
	return theme
}

type SiteThemeV1 struct {
	PrimaryRGB    string `json:"primary_rgb"`
	SecondaryRGB  string `json:"secondary_rgb"`
	BackgroundRGB string `json:"bg_rgb"`
	TextRGB       string `json:"text_rgb"`
	LogoSize      string `json:"logo_size"`
}
