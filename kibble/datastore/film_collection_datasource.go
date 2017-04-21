package datastore

import (
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

// Iterator - return a list of all films, iteration of 1
func (ds *FilmCollectionDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	clonedFilms := make([]*models.Film, len(ctx.Site.Films))
	for i, f := range ctx.Site.Films {
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
	return models.DataSourceError
}

// IsSlugMatch - is the slug a match
func (ds *FilmCollectionDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformFilm(f models.Film) *models.Film {
	f.Overview = models.ApplyContentTransforms(f.Overview)
	return &f
}
