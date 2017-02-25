package datastore

import (
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// FilmCollectionDataSource - a list of all films
type FilmCollectionDataSource struct{}

// GetName - returns the name of the datasource
func (ds *FilmCollectionDataSource) GetName() string {
	return "FilmCollection"
}

// Query - return the list of all films
func (ds *FilmCollectionDataSource) Query(req *http.Request) (jet.VarMap, error) {

	films, err := GetAllFilms()
	if err != nil || films == nil {
		return nil, err
	}

	clonedFilms := make([]*models.Film, len(*films))
	for i, f := range *films {
		clonedFilms[i] = transformFilm(f)
	}

	vars := make(jet.VarMap)
	vars.Set("films", clonedFilms)
	return vars, nil
}

// Iterator - return a list of all films, iteration of 1
func (ds *FilmCollectionDataSource) Iterator(route *models.Route, renderer models.Renderer) {

	films, _ := GetAllFilms()

	clonedFilms := make([]*models.Film, len(*films))
	for i, f := range *films {
		clonedFilms[i] = transformFilm(f)
	}

	vars := make(jet.VarMap)
	vars.Set("films", clonedFilms)
	renderer.Render(route, route.URLPath, vars)

}

func init() {
	AddDataSource(&FilmDataSource{})
	AddDataSource(&FilmCollectionDataSource{})
}

func transformFilm(f models.Film) *models.Film {
	f.Synopsis = applyContentTransforms(f.Synopsis)
	return &f
}

// GetAllFilms - returns all films
func GetAllFilms() (*[]models.Film, error) {
	return &allFilms, nil
}

// FindByID - find a film by its id
func FindByID(filmID int) (*models.Film, error) {
	for _, f := range allFilms {
		if f.ID == filmID {
			return &f, nil
		}
	}
	return nil, nil
}
