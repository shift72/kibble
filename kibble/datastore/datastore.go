package datastore

import (
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// DataSource - provides a set of data for querying and iterating over
type DataSource interface {
	GetName() string
	Query(*http.Request) (jet.VarMap, error)
	Iterator(route *models.Route, renderer models.Renderer)
	//ValidateRoute - check that the route contains valid tokens
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

func applyContentTransforms(data string) string {

	//TODO: apply shortcodes

	// apply mark down
	unsafe := blackfriday.MarkdownCommon([]byte(data))

	// return string(unsafe)
	// apply sanitization
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return string(html)
}
