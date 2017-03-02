package datastore

import (
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// FilmCollectionDataSource - a list of all films
type FilmCollectionDataSource struct{}

// GetName - returns the name of the datasource
func (ds *FilmCollectionDataSource) GetName() string {
	return "FilmCollection"
}

// GetEntityType - Get the entity type
func (ds *FilmCollectionDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Film{})
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
func (ds *FilmCollectionDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	films, _ := GetAllFilms()

	clonedFilms := make([]*models.Film, len(*films))
	for i, f := range *films {
		clonedFilms[i] = transformFilm(f)
	}

	vars := make(jet.VarMap)
	vars.Set("films", clonedFilms)
	vars.Set("site", ctx.Site)
	renderer.Render(ctx.Route, ctx.Route.URLPath, vars)

}

// GetRouteForEntity - get the route
func (ds *FilmCollectionDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *FilmCollectionDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return "!Error"
}

// IsSlugMatch - is the slug a match
func (ds *FilmCollectionDataSource) IsSlugMatch(slug string) bool {
	return false
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
