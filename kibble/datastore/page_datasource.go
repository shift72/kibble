package datastore

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// PageDataSource - single Page datasource
type PageDataSource struct{}

// GetName - name of the datasource
func (ds *PageDataSource) GetName() string {
	return "Page"
}

// GetEntityType - Get the entity type
func (ds *PageDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.Page{})
}

// Iterator - loop over each Page
func (ds *PageDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	data := make(jet.VarMap)

	for _, p := range ctx.Site.Pages {

		if p.PageType == "external" {
			continue // don't render external pages
		}

		data.Set("page", transformPage(p))
		data.Set("site", ctx.Site)
		renderer.Render(ctx.Route, ds.GetRouteForEntity(ctx, &p), data)
	}
}

// GetRouteForEntity - get the route
func (ds *PageDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Page)
	if ok {
		// special case for the home page
		if o.PageType == "homepage" {
			return ctx.RoutePrefix + "/"
		}
		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":slug", o.Slug, 1)
	}
	return models.DataSourceError
}

// GetRouteForSlug - get the route
func (ds *PageDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	//TODO: parse slug
	//TODO: fix errors
	p := strings.Split(slug, "/")
	pageID, _ := strconv.Atoi(p[2])
	page, _ := ctx.Site.Pages.FindPageByID(pageID)

	//TODO: does not work page, _ := ctx.Site.Pages.FindPageBySlug(slug)
	return ds.GetRouteForEntity(ctx, page)
}

// IsSlugMatch - checks if the slug is a match
func (ds *PageDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/page/")
}
