package models

import (
	"reflect"

	"github.com/CloudyKit/jet"
)

// Pagination describes a single page of results
type Pagination struct {
	Index       int
	Size        int
	Total       int
	PreviousURL string
	NextURL     string
}

// DataSource provides a set of data for querying and iterating over
type DataSource interface {
	GetName() string
	GetEntityType() reflect.Type
	Iterator(ctx RenderContext, renderer Renderer) (errorCount int)
	IsSlugMatch(slug string) bool
	GetRouteForEntity(ctx RenderContext, entity interface{}) string
	GetRouteForSlug(ctx RenderContext, slug string) string
	//TODO: ValidateRoute - check that the route contains valid tokens
}

// RenderContext - Context passed during rendering / serving
type RenderContext struct {
	Route       *Route
	RoutePrefix string
	Site        *Site
	Language    string
}

// Renderer - rendering implementation
type Renderer interface {
	Initialise()
	Render(route *Route, filePath string, data jet.VarMap) (errorCount int)
}

// Config - template configuration
type Config struct {
	DefaultLanguage string            `json:"defaultLanguage"`
	Languages       map[string]string `json:"languages"`
	Routes          []Route           `json:"routes"`
	SiteURL         string            `json:"siteUrl"`
	BuilderVersion  string            `json:"builderVersion"`
	Version         string            `json:"version"`
	Name            string            `json:"Name"`
	Private         PrivateConfig     `json:"-"`
	DisableCache    bool              `json:"-"`
	RunAsAdmin      bool              `json:"-"`
	SkipLogin       bool              `json:"-"`
}

// PrivateConfig - config loaded from
type PrivateConfig struct {
	APIKey string `json:"apikey"`
}
