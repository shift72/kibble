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

func loadSiteBrand(cfg *models.Config, site *models.Site) error {
	//If both self_service toggles are off
	if !site.Toggles["self_service_site_images"] && !site.Toggles["self_service_css"] {
		// Do nothing - use local site assets
		return nil
	}

	//Load from API
	path := fmt.Sprintf("%s/services/users/v1/site_brand", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return fmt.Errorf(" site branding from API failed to load %s", err)
	}

	var siteBrand SiteBrandsV1

	err = json.Unmarshal([]byte(data), &siteBrand)
	if err != nil {
		return err
	}

	for i, Info := range siteBrand.Images {
		println(Info.Type)
		println(Info.URL)
		println(siteBrand.Images[i].Type)
		println(siteBrand.Images[i].URL)
	}
	var strings = siteBrand.Images[0]
	print(strings)
	// site.SiteBrand.Links = siteBrand.Links

	return nil
}

type SiteBrandV1 struct {
	Images []map[string]string `json:"images,omitempty"`
	Links  []map[string]string `json:"links,omitempty"`
}
type SiteBrandsV1 struct {
	Images []InfoV1 `json:"images,omitempty"`
	Links  []InfoV1 `json:"links,omitempty"`
}
type InfoV1 struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}
