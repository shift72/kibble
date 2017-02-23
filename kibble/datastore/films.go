package datastore

import "github.com/indiereign/shift72-kibble/kibble/models"

//TODO: register the datasources here

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
