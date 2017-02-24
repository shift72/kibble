package datastore

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/pressly/chi"
)

// dummy data
var allFilms = []models.Film{
	{
		ID:    1,
		Title: "Forrest Gump",
	},
	{
		ID:    2,
		Title: "Angel at my table",
	},
}

var film = &DataSource{
	Name: "Film",
	Query: func(req *http.Request) (jet.VarMap, error) {

		filmID, err := strconv.Atoi(chi.URLParam(req, "filmID"))
		if err != nil {
			return nil, err
		}

		f, err := FindByID(filmID)
		if err != nil || f == nil {
			return nil, err
		}

		vars := make(jet.VarMap)
		vars.Set("film", f)
		return vars, nil

	},
	Iterator: func(route *models.Route, yield func(route *models.Route, path string, data jet.VarMap)) {

		films, _ := GetAllFilms()
		data := make(jet.VarMap)
		for _, f := range *films {

			p := strings.Replace(route.URLPath, ":filmID", strconv.Itoa(f.ID), 1)
			data.Set("film", f)
			yield(route, p, data)

		}
	},
}

var filmCollection = &DataSource{
	Name: "FilmCollection",
	Query: func(req *http.Request) (jet.VarMap, error) {

		f, err := GetAllFilms()
		if err != nil || f == nil {
			return nil, err
		}

		vars := make(jet.VarMap)
		vars.Set("films", f)
		return vars, nil

	},
	Iterator: func(route *models.Route, yield func(*models.Route, string, jet.VarMap)) {

		films, _ := GetAllFilms()
		vars := make(jet.VarMap)
		vars.Set("films", films)
		yield(route, route.URLPath, vars)

	},
}

func init() {
	AddDataSource(film)
	AddDataSource(filmCollection)
}

// GetAllFilms - returns all films
func GetAllFilms() (*[]models.Film, error) {
	return &allFilms, nil
}

// FindByID - dd
func FindByID(filmID int) (*models.Film, error) {
	for _, f := range allFilms {
		if f.ID == filmID {
			return &f, nil
		}
	}
	return nil, nil
}
