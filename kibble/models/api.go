package models

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

// PageFeature - part of a page
type PageFeature struct {
	FeatureID   int      `json:"feature_id"`
	Layout      string   `json:"layout"`
	ItemsPerRow int      `json:"items_per_row"`
	ItemLayout  string   `json:"item_layout"`
	Slug        string   `json:"slug"`
	DisplayName *string  `json:"display_name"`
	Items       []string `json:"items"`
	// ResolvedItems?       []interface `json:"-"`
}

// Page - page structure
type Page struct {
	CarouselImage  *string       `json:"carousel_image"`
	Content        string        `json:"content"`
	HeaderImage    *string       `json:"header_image"`
	ID             int           `json:"id"`
	LandscapeImage *string       `json:"landscape_image"`
	PageFeatures   []interface{} `json:"page_features"`
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
