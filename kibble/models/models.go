package models

import (
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet"
)

// Film - represents a film
type Film struct {
	ID       int
	Slug     string
	Title    string
	Synopsis string
}

// Route - represents a route for rendering and
type Route struct {
	Name               string       `json:"name"`
	URLPath            string       `json:"urlPath"`
	TemplatePath       string       `json:"templatePath"`
	DataSource         string       `json:"datasource"`
	ResolvedDataSouce  DataSource   `json:"-"`
	ResolvedEntityType reflect.Type `json:"-"`
}

// RouteRegistry - stores a list of routes
type RouteRegistry struct {
	routes []*Route
}

// DataSource - provides a set of data for querying and iterating over
type DataSource interface {
	GetName() string
	GetEntityType() reflect.Type
	Query(*http.Request) (jet.VarMap, error)
	Iterator(ctx RenderContext, renderer Renderer)
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
}

// Renderer - rendering implementation
type Renderer interface {
	Render(route *Route, filePath string, data jet.VarMap)
}

// Config -
type Config struct {
	DefaultLanguage string            `json:"defaultLanguage"`
	Languages       map[string]string `json:"languages"`
	Routes          []Route           `json:"routes"`
	SiteURL         string            `json:"siteUrl"`
}
