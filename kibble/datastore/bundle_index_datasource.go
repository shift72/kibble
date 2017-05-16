package datastore

import (
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// BundleIndexDataSource - a list of all bundles
type BundleIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *BundleIndexDataSource) GetName() string {
	return "BundleIndex"
}

// GetEntityType - Get the entity type
func (ds *BundleIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Bundle{})
}

// Iterator - return a list of all bundles, iteration of 1
func (ds *BundleIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) int {

	clonedBundles := make([]*models.Bundle, len(ctx.Site.Bundles))
	for i, f := range ctx.Site.Bundles {
		clonedBundles[i] = transformBundle(f)
	}

	vars := make(jet.VarMap)
	vars.Set("bundles", clonedBundles)
	vars.Set("site", ctx.Site)
	return renderer.Render(ctx.Route, ctx.RoutePrefix+ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *BundleIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *BundleIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *BundleIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformBundle(f models.Bundle) *models.Bundle {
	f.Description = models.ApplyContentTransforms(f.Description)
	return &f
}
