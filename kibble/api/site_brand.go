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

func LoadSiteBrand(cfg *models.Config, site *models.Site) error {
	//If both self_service toggles are off
	if !site.Toggles["self_service_site_images"] && !site.Toggles["self_service_css"] {
		// Do nothing - use local site assets
		log.Infof("Self Service Images and CSS disabled, all branding assets will use default paths")
		return nil
	}

	//Load from API
	path := fmt.Sprintf("%s/services/users/v1/site_brand", cfg.SiteURL)

	data, err := Get(cfg, path)
	if err != nil {
		return fmt.Errorf(" Site branding from API failed to load %s", err)
	}

	var siteBrands SiteBrandsV1

	err = json.Unmarshal([]byte(data), &siteBrands)
	if err != nil {
		return err
	}

	if site.Toggles["self_service_site_images"] {

		site.SiteBrand.Images = mapBranding(siteBrands.Images)
		log.Infof("Self Service Images enabled:")
		printAssets(site.SiteBrand.Images)

	}
	if site.Toggles["self_service_css"] {
		site.SiteBrand.Links = mapBranding(siteBrands.Links)
		log.Infof("Self Service Links enabled:")
		printAssets(site.SiteBrand.Links)
	}
	return nil
}

func mapBranding(brandingItems []SiteBrandItemV1) map[string]string {
	assetMap := make(map[string]string)
	for _, asset := range brandingItems {
		assetMap[asset.Type] = asset.URL
	}
	return assetMap
}

func printAssets(assetMap map[string]string) {
	for s := range assetMap {
		log.Infof(" %s", s)
	}
}

type SiteBrandsV1 struct {
	Images []SiteBrandItemV1 `json:"images,omitempty"`
	Links  []SiteBrandItemV1 `json:"links,omitempty"`
}
type SiteBrandItemV1 struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}
