package models

import (
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet"
)

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
}

var store map[string]DataSource

// AddDataSource - register a datasource
func AddDataSource(ds DataSource) {

	if store == nil {
		store = make(map[string]DataSource)
	}

	store[ds.GetName()] = ds
}

// FindDataSource - find the data source by name
func FindDataSource(name string) DataSource {
	return store[name]
}
