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

// FeatureToggles - store feature toggles
type FeatureToggles map[string]bool

// Site -
type Site struct {
	Config      ServiceConfig
	SiteConfig  *Config
	Toggles     FeatureToggles
	Navigation  Navigation
	Languages   []Language
	Pages       Pages
	Films       FilmCollection
	TVShows     TVShowCollection
	TVSeasons   TVSeasonCollection
	TVEpisodes  TVEpisodeCollection
	Bundles     BundleCollection
	Collections CollectionCollection
	Plans       PlanCollection
	Taxonomies  Taxonomies
}

// Language - instance of a language
type Language struct {
	Code               string
	Locale             string
	DefinitionFilePath string
	IsDefault          bool
}

// ImageSet - set of images
type ImageSet struct {
	Background     string
	Carousel       string
	Landscape      string
	Portrait       string
	Header         string
	Classification string
	Seo            string
}

// Seo - common seo settings
type Seo struct {
	SiteName    string
	Title       string
	Description string
	Keywords    string
	Image       string
	VideoURL    string
}

// NavigationItem - nestable structure
type NavigationItem struct {
	Label string `json:"label"`
	Link  struct {
		PageID      int    `json:"page_id"`
		Slug        string `json:"slug"`
		ExternalURL string `json:"url"`
	} `json:"link"`
	Items []NavigationItem `json:"items"`
}

// Navigation - header and footer
type Navigation struct {
	Footer []NavigationItem `json:"footer"`
	Header []NavigationItem `json:"header"`
}

// Trailer -
type Trailer struct {
	URL  string
	Type string
}

// CastMember -
type CastMember struct {
	Name      string
	Character string
}

// CrewMember -
type CrewMember struct {
	Name string
	Job  string
}

// SubtitleTrack -
type SubtitleTrack struct {
	Language string
	Name     string
	Type     string
	Path     string
}

// Classification -
type Classification struct {
	CountryCode string
	Label       string
	Description string
}

// CustomFields are key-value pairs that can be aded to a film, season, bonus, or episode
type CustomFields map[string]interface{}

// GetString returns the custom field in string format
func (fields CustomFields) GetString(fieldKey string, defaultValue string) string {
	if value, ok := fields[fieldKey]; ok {
		if castValue, ok := value.(string); ok {
			return castValue
		}
	}

	return defaultValue
}

// GetBool returns the custom field in bool format
func (fields CustomFields) GetBool(fieldKey string, defaultValue bool) bool {
	if value, ok := fields[fieldKey]; ok {
		if castValue, ok := value.(bool); ok {
			return castValue
		}
	}

	return defaultValue
}

// GetNumber returns the custom field in float64 format
func (fields CustomFields) GetNumber(fieldKey string, defaultValue float64) float64 {
	if value, ok := fields[fieldKey]; ok {
		if castValue, ok := value.(float64); ok {
			return castValue
		}
	}

	return defaultValue
}
