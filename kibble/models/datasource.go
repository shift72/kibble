package models

import "errors"

var (
	// ErrDataSource is returned when no datasource could be found
	ErrDataSource = "!Error"
	// ErrDataSourceMissing raises when the datasource could not be found
	ErrDataSourceMissing = errors.New("Missing")
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
