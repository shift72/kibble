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
	"strings"
)

func loadSiteBrand(cfg *models.Config, site *models.Site) error {
	//If both self_service toggles are off
	if !site.Toggles["self_service_site_images"] && !site.Toggles["self_service_css"] {
		// Do nothing - use local site assets
		log.Infof("Self Service Images and CSS disabled, using local branding assets")
		return nil
	}

	//Load from API
	path := fmt.Sprintf("%s/services/users/v1/site_brand", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return fmt.Errorf(" Site branding from API failed to load %s", err)
	}

	var siteBrand SiteBrandsV1

	err = json.Unmarshal([]byte(data), &siteBrand)
	if err != nil {
		return err
	}

	//initialise empty maps, empty is still valid on site model
	images := make(map[string]string)
	links := make(map[string]string)

	var brandingState strings.Builder
	brandingState.WriteString("Branding Assets:")

	if site.Toggles["self_service_site_images"] {
		log.Info("Self Service Images Enbabled")
		brandingState.WriteString("\nImages: ")
		for _, Info := range siteBrand.Images {
			images[Info.Type] = Info.URL
			brandingState.WriteString(Info.Type + " ")
		}
	}
	site.SiteBrand.Images = images

	if site.Toggles["self_service_css"] {
		log.Info("Self Service Links Enbabled")
		brandingState.WriteString("\nLinks: ")
		for _, Link := range siteBrand.Links {
			links[Link.Type] = Link.URL
			brandingState.WriteString(Link.Type + " ")
		}
	}
	site.SiteBrand.Links = links

	log.Info("%s", brandingState.String())

	return nil
}

type SiteBrandsV1 struct {
	Images []SiteBrandItemV1 `json:"images,omitempty"`
	Links  []SiteBrandItemV1 `json:"links,omitempty"`
}
type SiteBrandItemV1 struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}
