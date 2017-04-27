package datastore

import (
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// Init -
func Init() {
	models.AddDataSource(&FilmDataSource{})
	models.AddDataSource(&FilmIndexDataSource{})

	models.AddDataSource(&PageDataSource{})
	models.AddDataSource(&PageIndexDataSource{})

	models.AddDataSource(&BundleDataSource{})
	models.AddDataSource(&BundleIndexDataSource{})
}
