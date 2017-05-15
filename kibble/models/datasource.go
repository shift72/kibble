package models

import "errors"

var (
	// DataSourceError - Error to be written to the template
	DataSourceError   = "!Error"
	DataSourceMissing = errors.New("Missing")
)

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
