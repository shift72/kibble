package datastore

import (
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// PageCollectionDataSource - a list of all Pages
type PageCollectionDataSource struct{}

// GetName - returns the name of the datasource
func (ds *PageCollectionDataSource) GetName() string {
	return "PageCollection"
}

// GetEntityType - Get the entity type
func (ds *PageCollectionDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Page{})
}

// Query - return the list of all Pages
func (ds *PageCollectionDataSource) Query(ctx models.RenderContext, req *http.Request) (jet.VarMap, error) {

	clonedPages := make([]*models.Page, len(ctx.Site.Pages))
	for i, f := range ctx.Site.Pages {
		clonedPages[i] = transformPage(f)
	}

	vars := make(jet.VarMap)
	vars.Set("pages", clonedPages)
	vars.Set("site", ctx.Site)
	return vars, nil
}

// Iterator - return a list of all Pages, iteration of 1
func (ds *PageCollectionDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	clonedPages := make([]*models.Page, len(ctx.Site.Pages))
	for i, f := range ctx.Site.Pages {
		clonedPages[i] = transformPage(f)
	}

	vars := make(jet.VarMap)
	vars.Set("pages", clonedPages)
	vars.Set("site", ctx.Site)
	renderer.Render(ctx.Route, ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *PageCollectionDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *PageCollectionDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.DataSourceError
}

// IsSlugMatch - is the slug a match
func (ds *PageCollectionDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformPage(f models.Page) *models.Page {
	f.Content = applyContentTransforms(f.Content)
	return &f
}
