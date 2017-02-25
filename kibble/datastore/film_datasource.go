package datastore

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/pressly/chi"
)

// FilmDataSource - single film datasource
type FilmDataSource struct{}

// GetName - name of the datasource
func (ds *FilmDataSource) GetName() string {
	return "Film"
}

// Query - return a single film
func (ds *FilmDataSource) Query(req *http.Request) (jet.VarMap, error) {

	filmID, err := strconv.Atoi(chi.URLParam(req, "filmID"))
	if err != nil {
		return nil, err
	}

	f, err := FindByID(filmID)
	if err != nil || f == nil {
		return nil, err
	}
	c := transformFilm(*f)

	vars := make(jet.VarMap)
	vars.Set("film", c)
	return vars, nil
}

// Iterator - loop over each film
func (ds *FilmDataSource) Iterator(route *models.Route, renderer models.Renderer) {

	films, _ := GetAllFilms()
	data := make(jet.VarMap)

	for _, f := range *films {

		filePath := strings.Replace(route.URLPath, ":filmID", strconv.Itoa(f.ID), 1)

		c := transformFilm(f)

		data.Set("film", c)
		renderer.Render(route, filePath, data)
	}
}
