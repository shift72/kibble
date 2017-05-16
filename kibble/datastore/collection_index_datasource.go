package datastore

import (
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// CollectionIndexDataSource - a list of all Collections
type CollectionIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *CollectionIndexDataSource) GetName() string {
	return "CollectionIndex"
}

// GetEntityType - Get the entity type
func (ds *CollectionIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Collection{})
}

// Iterator - return a list of all Collections, iteration of 1
func (ds *CollectionIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	clonedCollections := make([]*models.Collection, len(ctx.Site.Collections))
	for i, f := range ctx.Site.Collections {
		clonedCollections[i] = transformCollection(f)
	}

	vars := make(jet.VarMap)
	vars.Set("collections", clonedCollections)
	vars.Set("site", ctx.Site)
	return renderer.Render(ctx.Route, ctx.RoutePrefix+ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *CollectionIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *CollectionIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *CollectionIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformCollection(f models.Collection) *models.Collection {
	f.Description = models.ApplyContentTransforms(f.Description)
	return &f
}
