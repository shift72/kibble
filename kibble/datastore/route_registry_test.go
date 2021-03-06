package datastore

import (
	"kibble/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFilmDataSourceRoutesFromConfig(t *testing.T) {

	models.AddDataSource(&FilmDataSource{})

	cfg := models.Config{
		Routes: []models.Route{
			{
				Name:         "filmItem",
				URLPath:      "/film-special/:filmID",
				TemplatePath: "film/item.jet",
				DataSource:   "Film",
			},
		},
	}

	routeRegistry, err := models.NewRouteRegistryFromConfig(&cfg)
	assert.NoError(t, err)

	routes := routeRegistry.GetAll()

	assert.Equal(t, routes[0].TemplatePath, "film/item.jet")
}

func TestLoadDefaultFileSystemRouteFromConfig(t *testing.T) {

	models.AddDataSource(&FileSystemDataSource{})

	cfg := models.Config{
		Routes: []models.Route{},
	}

	routeRegistry, err := models.NewRouteRegistryFromConfig(&cfg)
	assert.NoError(t, err)

	routes := routeRegistry.GetAll()

	assert.Equal(t, routes[0].TemplatePath, ".")
}
