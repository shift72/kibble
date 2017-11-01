package datastore

import (
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// TVSeasonIndexDataSource - a list of all films
type TVSeasonIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *TVSeasonIndexDataSource) GetName() string {
	return "TVSeasonIndex"
}

// GetEntityType - Get the entity type
func (ds *TVSeasonIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.TVSeason{})
}

// Iterator - return a list of all tv seasons, iteration of 1
func (ds *TVSeasonIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	cloned := make([]*models.TVSeason, len(ctx.Site.TVSeasons))
	for i, s := range ctx.Site.TVSeasons {
		cloned[i] = transformTVSeason(s)
	}

	vars := make(jet.VarMap)
	vars.Set("tvseasons", cloned)
	vars.Set("site", ctx.Site)
	return renderer.Render(ctx.Route, ctx.RoutePrefix+ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *TVSeasonIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *TVSeasonIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *TVSeasonIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformTVSeason(f models.TVSeason) *models.TVSeason {
	f.Overview = models.ApplyContentTransforms(f.Overview)
	return &f
}
