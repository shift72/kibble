package datastore

import (
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// BundleCollectionDataSource - a list of all bundles
type BundleCollectionDataSource struct{}

// GetName - returns the name of the datasource
func (ds *BundleCollectionDataSource) GetName() string {
	return "BundleCollection"
}

// GetEntityType - Get the entity type
func (ds *BundleCollectionDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Bundle{})
}

// Query - return the list of all bundles
func (ds *BundleCollectionDataSource) Query(ctx models.RenderContext, req *http.Request) (jet.VarMap, error) {

	clonedBundles := make([]*models.Bundle, len(ctx.Site.Bundles))
	for i, f := range ctx.Site.Bundles {
		clonedBundles[i] = transformBundle(f)
	}

	vars := make(jet.VarMap)
	vars.Set("bundles", clonedBundles)
	vars.Set("site", ctx.Site)
	return vars, nil
}

// Iterator - return a list of all bundles, iteration of 1
func (ds *BundleCollectionDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	clonedBundles := make([]*models.Bundle, len(ctx.Site.Bundles))
	for i, f := range ctx.Site.Bundles {
		clonedBundles[i] = transformBundle(f)
	}

	vars := make(jet.VarMap)
	vars.Set("bundles", clonedBundles)
	vars.Set("site", ctx.Site)
	renderer.Render(ctx.Route, ctx.Route.URLPath, vars)

}

// GetRouteForEntity - get the route
func (ds *BundleCollectionDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *BundleCollectionDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.DataSourceError
}

// IsSlugMatch - is the slug a match
func (ds *BundleCollectionDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformBundle(f models.Bundle) *models.Bundle {
	f.Description = applyContentTransforms(f.Description)
	return &f
}
