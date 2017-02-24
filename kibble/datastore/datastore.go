package datastore

import (
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// DataSource - provides a set of data for querying and iterating over
type DataSource struct {
	Name     string
	Query    func(*http.Request) (jet.VarMap, error)
	Iterator func(route *models.Route, yeild func(*models.Route, string, jet.VarMap))
	//ValidateRoute - check that the route contains valid tokens
}

var store map[string]*DataSource

// AddDataSource - register a datasource
func AddDataSource(ds *DataSource) {

	if store == nil {
		store = make(map[string]*DataSource)
	}

	store[ds.Name] = ds
}

// FindDataSource - find the data source by name
func FindDataSource(name string) *DataSource {
	return store[name]
}
