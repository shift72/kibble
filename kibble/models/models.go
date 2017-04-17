package models

import (
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/pressly/chi"
)

// Pagination -
type Pagination struct {
	Index       int
	Size        int
	Total       int
	PreviousURL string
	NextURL     string
}

// DataSource - provides a set of data for querying and iterating over
type DataSource interface {
	GetName() string
	GetEntityType() reflect.Type
	Query(ctx RenderContext, req *http.Request) (jet.VarMap, error)
	Iterator(ctx RenderContext, renderer Renderer)
	IsSlugMatch(slug string) bool
	GetRouteForEntity(ctx RenderContext, entity interface{}) string
	GetRouteForSlug(ctx RenderContext, slug string) string
	RegisterRoutes(router chi.Router, route *Route, handler func(w http.ResponseWriter, req *http.Request))
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
	Render(route *Route, filePath string, data jet.VarMap)
}

// Config -
type Config struct {
	DefaultLanguage string            `json:"defaultLanguage"`
	Languages       map[string]string `json:"languages"`
	Routes          []Route           `json:"routes"`
	SiteURL         string            `json:"siteUrl"`
	Private         PrivateConfig     `json:"-"`
}

// PrivateConfig - config loaded from
type PrivateConfig struct {
	APIKey string `json:"apikey"`
}
