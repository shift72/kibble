//    Copyright 2018 SHIFT72
//
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

package models

import (
	"strconv"

	"github.com/indiereign/shift72-kibble/kibble/utils"
)

// ServiceConfig -
type ServiceConfig map[string]string

// SelectDefaultImageType - select the default image
func (cfg ServiceConfig) SelectDefaultImageType(landscape, portrait string) string {
	if cfg["default_image_type"] == "landscape" {
		return landscape
	}
	return portrait
}

// ForceAbsoluteImagePath fixes up relative image path by prefixing the `image_root_path` to it.
func (cfg ServiceConfig) ForceAbsoluteImagePath(url string) string {
	if len(url) > 0 {
		return cfg["image_root_path"] + url
	}

	return url
}

// GetSEOTitle - get the seo title
func (cfg ServiceConfig) GetSEOTitle(seoTitle, title string) string {
	return utils.Join(" ",
		cfg["seo_title_prefix"],
		utils.Coalesce(seoTitle, title),
		cfg["seo_title_suffix"])
}

// GetSiteName - get the site name
func (cfg ServiceConfig) GetSiteName() string {
	return cfg["seo_site_name"]
}

// GetKeywords - get the keywords, appending any passed keywords
func (cfg ServiceConfig) GetKeywords(keywords string) string {
	return utils.Join(", ", cfg["seo_site_keywords"], keywords)
}

// GetInt - get the key and cast
func (cfg ServiceConfig) GetInt(key string, defaultArgs ...int) int {

	d := 0
	if len(defaultArgs) > 0 {
		d = defaultArgs[0]
	}

	s, ok := cfg[key]
	if ok {
		i, err := strconv.Atoi(s)
		if err != nil {
			return d
		}
		return i
	}
	return d
}
