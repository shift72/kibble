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
	"time"
)

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

// GenericItem - used to store the common properties
type GenericItem struct {
	// link to the actual item
	InnerItem interface{}
	// film / show / season/ episode / bundle / page
	ItemType string
	Slug     string
	Title    string
	Images   ImageSet
	Seo      Seo
}

// PageCollection - part of a page
type PageCollection struct {
	ID          int
	Layout      string
	ItemsPerRow int
	ItemLayout  string
	Slug        string
	TitleSlug   string
	DisplayName string
	Items       []GenericItem
}

// Page - page structure
type Page struct {
	ID              int
	Slug            string
	Title           string
	TitleSlug       string
	Content         string
	Tagline         string
	Seo             Seo
	Images          ImageSet
	PageCollections []PageCollection
	PageType        string
	URL             string
}

// Pages -
type Pages []Page

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

// FilmBonusCollection - all films
type FilmBonusCollection []FilmBonus

// FilmBonus - film bonus model
type FilmBonus struct {
	Slug      string
	Number    int
	Title     string
	Images    ImageSet
	Subtitles []SubtitleTrack
	Runtime   Runtime
	Overview  string
}

// FilmCollection - all films
type FilmCollection []Film

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

// Film - all of the film bits
type Film struct {
	ID              int
	Slug            string
	Title           string
	TitleSlug       string
	Trailers        []Trailer
	Bonuses         FilmBonusCollection
	Cast            []CastMember
	Crew            []CrewMember
	Studio          []string
	Overview        string
	Tagline         string
	ReleaseDate     time.Time
	Runtime         Runtime
	Countries       StringCollection
	Languages       StringCollection
	Genres          StringCollection
	Seo             Seo
	Images          ImageSet
	Recommendations []GenericItem
	Subtitles       []SubtitleTrack
}

// TVShow -
type TVShow struct {
	ID               int
	Slug             string
	Trailers         []Trailer
	Genres           StringCollection
	Overview         string
	Countries        StringCollection
	Languages        StringCollection
	ReleaseDate      time.Time
	Tagline          string
	Studio           []string
	Title            string
	TitleSlug        string
	AvailableSeasons []string
	Seasons          TVSeasonCollection
	Images           ImageSet
}

// TVEpisode -
type TVEpisode struct {
	Title         string
	EpisodeNumber int
	Overview      string
	Runtime       Runtime
	Images        ImageSet
	Subtitles     []SubtitleTrack
}

// TVSeason -
type TVSeason struct {
	Slug            string
	SeasonNumber    int
	Title           string
	Tagline         string
	Overview        string
	PublishingState string
	ShowInfo        *TVShow
	Seo             Seo
	Images          ImageSet
	Trailers        []Trailer
	Episodes        []TVEpisode
	Cast            []CastMember
	Crew            []CrewMember
	Recommendations []GenericItem
}

// StringCollection - Allows us to add methods to []string for easing UI array usage
type StringCollection []string

// Runtime - Allows us to get accurate measures of hours and minutes.
type Runtime int

// TVShowCollection -
type TVShowCollection []TVShow

// TVSeasonCollection -
type TVSeasonCollection []TVSeason
