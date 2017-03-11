package datastore

import (
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// Init -
func Init() {
	models.AddDataSource(&FilmDataSource{})
	models.AddDataSource(&FilmCollectionDataSource{})

	models.AddDataSource(&PageDataSource{})
	models.AddDataSource(&PageCollectionDataSource{})

	models.AddDataSource(&BundleDataSource{})
	models.AddDataSource(&BundleCollectionDataSource{})
}
