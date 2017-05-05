package datastore

import (
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// FilmIndexDataSource - a list of all films
type FilmIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *FilmIndexDataSource) GetName() string {
	return "FilmIndex"
}

// GetEntityType - Get the entity type
func (ds *FilmIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Film{})
}

// Iterator - return a list of all films, iteration of 1
func (ds *FilmIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	clonedFilms := make([]*models.Film, len(ctx.Site.Films))
	for i, f := range ctx.Site.Films {
		clonedFilms[i] = transformFilm(f)
	}

	vars := make(jet.VarMap)
	vars.Set("films", clonedFilms)
	vars.Set("site", ctx.Site)
	renderer.Render(ctx.Route, ctx.RoutePrefix+ctx.Route.URLPath, vars)

}

// GetRouteForEntity - get the route
func (ds *FilmIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *FilmIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.DataSourceError
}

// IsSlugMatch - is the slug a match
func (ds *FilmIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformFilm(f models.Film) *models.Film {
	f.Overview = models.ApplyContentTransforms(f.Overview)
	return &f
}