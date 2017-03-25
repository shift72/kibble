package models

import (
	"encoding/json"
	"time"
)

// ServiceConfig -
type ServiceConfig map[string]string

// FeatureToggles - store feature toggles
type FeatureToggles map[string]bool

// Site -
type Site struct {
	Config     ServiceConfig
	Toggles    FeatureToggles
	Navigation Navigation
	Pages      PageCollection
	Films      FilmCollection
	Bundles    BundleCollection
	Taxonomies Taxonomies
}

// "page_features": [{
//     "feature_id": 125,
//     "layout": "slider",
//     "items_per_row": 3,
//     "item_layout": "portrait",
//     "slug": "test-01234",
//     "display_name": null,
//     "items": ["/film/121"]
// },

// ImageSet - set of images
type ImageSet struct {
	BackgroundImage *string
	CarouselImage   *string
	LandscapeImage  *string
	PortraitImage   *string
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
}

// PageFeature - part of a page
type PageFeature struct {
	FeatureID     int           `json:"feature_id"`
	Layout        string        `json:"layout"`
	ItemsPerRow   int           `json:"items_per_row"`
	ItemLayout    string        `json:"item_layout"`
	Slug          string        `json:"slug"`
	DisplayName   *string       `json:"display_name"`
	Items         []string      `json:"items"`
	ResolvedItems []GenericItem `json:"-"`
}

// Page - page structure
type Page struct {
	CarouselImage  *string       `json:"carousel_image"`
	Content        string        `json:"content"`
	HeaderImage    *string       `json:"header_image"`
	ID             int           `json:"id"`
	LandscapeImage *string       `json:"landscape_image"`
	PageFeatures   []PageFeature `json:"page_features"`
	PageType       string        `json:"page_type"`
	PortraitImage  *string       `json:"portrait_image"`
	SeoDescription *string       `json:"seo_description"`
	SeoKeywords    *string       `json:"seo_keywords"`
	SeoTitle       *string       `json:"seo_title"`
	Slug           string        `json:"slug"`
	Tagline        *string       `json:"tagline"`
	Title          string        `json:"title"`
	URL            string        `json:"url"`
}

// PageCollection -
type PageCollection []Page

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

// Bios - contains all pages and navigation items
type Bios struct {
	Navigation Navigation     `json:"navigation"`
	Pages      PageCollection `json:"pages"`
}

// FilmSummary - summary of the Film
type FilmSummary struct {
	BackgroundImage     interface{} `json:"background_image"`
	CarouselImage       interface{} `json:"carousel_image"`
	ClassificationImage interface{} `json:"classification_image"`
	HeaderImage         interface{} `json:"header_image"`
	ID                  int         `json:"id"`
	ImdbID              interface{} `json:"imdb_id"`
	LandscapeImage      interface{} `json:"landscape_image"`
	PortraitImage       string      `json:"portrait_image"`
	PublishedDate       string      `json:"published_date"`
	Slug                string      `json:"slug"`
	StatusID            int         `json:"status_id"`
	Title               string      `json:"title"`
}

// FilmBonusCollection - all films
type FilmBonusCollection []FilmBonus

// FilmBonus - film bonus model
type FilmBonus struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	ImageUrls struct {
		Portrait       string `json:"portrait"`
		Landscape      string `json:"landscape"`
		Header         string `json:"header"`
		Carousel       string `json:"carousel"`
		Bg             string `json:"bg"`
		Classification string `json:"classification"`
	} `json:"image_urls"`
	SubtitleTracks []interface{} `json:"subtitle_tracks"`
}

// FilmCollection - all films
type FilmCollection []Film

// Film - all of the film bits
type Film struct {
	Trailers []struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"trailers,omitempty"`
	//TODO: add a bonus struct
	Bonuses FilmBonusCollection `json:"bonuses"`
	Cast    []struct {
		Name      string `json:"name"`
		Character string `json:"character"`
	} `json:"cast"`
	Crew []struct {
		Name string `json:"name"`
		Job  string `json:"job"`
	} `json:"crew"`
	Studio []struct {
		Name string `json:"name"`
	} `json:"studio"`
	Overview    string    `json:"overview"`
	Tagline     string    `json:"tagline"`
	ReleaseDate time.Time `json:"release_date"`
	Runtime     float32   `json:"runtime"`
	Countries   []string  `json:"countries"`
	Languages   []string  `json:"languages"`
	Genres      []string  `json:"genres"`
	Title       string    `json:"title"`
	TitleSlug   string    `json:"-"`
	Slug        string    `json:"slug"`
	FilmID      int       `json:"film_id"`
	ID          int       `json:"id"`
	ImageUrls   struct {
		Portrait       string `json:"portrait"`
		Landscape      string `json:"landscape"`
		Header         string `json:"header"`
		Carousel       string `json:"carousel"`
		Bg             string `json:"bg"`
		Classification string `json:"classification"`
	} `json:"image_urls"`
	Recommendations         []string      `json:"recommendations"`
	ResolvedRecommendations []GenericItem `json:"-"`
	//TODO: add a subtitle tracks struct
	SubtitleTracks []interface{} `json:"subtitle_tracks"`
	Subtitles      []string      `json:"-"`
	// manage the inconsistent api
	SubtitlesRaw json.RawMessage `json:"subtitles,omitempty"`
}

// BundleCollection - all bundles
type BundleCollection []Bundle

// Bundle - model
type Bundle struct {
	ID             int           `json:"id"`
	Slug           string        `json:"-"`
	Title          string        `json:"title"`
	TitleSlug      string        `json:"-"`
	Tagline        string        `json:"tagline"`
	Description    string        `json:"description"`
	Status         string        `json:"status"`
	PublishedDate  time.Time     `json:"published_date"`
	SeoTitle       string        `json:"seo_title"`
	SeoKeywords    string        `json:"seo_keywords"`
	SeoDescription string        `json:"seo_description"`
	PortraitImage  string        `json:"portrait_image"`
	LandscapeImage string        `json:"landscape_image"`
	CarouselImage  string        `json:"carousel_image"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	BgImage        string        `json:"bg_image"`
	PromoURL       string        `json:"promo_url"`
	ExternalID     interface{}   `json:"external_id"`
	Items          []string      `json:"items"`
	ResolvedItems  []GenericItem `json:"-"`
}
