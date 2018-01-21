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

	models.AddDataSource(&CollectionDataSource{})
	models.AddDataSource(&CollectionIndexDataSource{})

	models.AddDataSource(&TVShowDataSource{})
	models.AddDataSource(&TVShowIndexDataSource{})

	models.AddDataSource(&TVSeasonDataSource{})
	models.AddDataSource(&TVSeasonIndexDataSource{})
}
