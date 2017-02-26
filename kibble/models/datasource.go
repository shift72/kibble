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
	Iterator(route *Route, renderer Renderer)
	IsSlugMatch(slug string) bool
	GetRouteForEntity(route *Route, entity interface{}) string
	GetRouteForSlug(route *Route, slug string) string
	//TODO: ValidateRoute - check that the route contains valid tokens
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
