package datastore

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/pressly/chi"
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

// Query - return a single Page
func (ds *PageDataSource) Query(ctx models.RenderContext, req *http.Request) (jet.VarMap, error) {

	pageID, err := strconv.Atoi(chi.URLParam(req, "pageID"))
	if err != nil {
		return nil, err
	}

	p, err := ctx.Site.Pages.FindPageByID(pageID)
	if err != nil || p == nil {
		return nil, err
	}
	c := transformPage(*p)

	vars := make(jet.VarMap)
	vars.Set("page", c)
	vars.Set("site", ctx.Site)
	return vars, nil
}

// Iterator - loop over each Page
func (ds *PageDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	data := make(jet.VarMap)

	for _, p := range ctx.Site.Pages {

		filePath := ds.GetRouteForEntity(ctx, &p)

		c := transformPage(p)

		data.Set("page", c)
		data.Set("site", ctx.Site)
		renderer.Render(ctx.Route, filePath, data)
	}
}

// GetRouteForEntity - get the route
func (ds *PageDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Page)
	if ok {
		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":pageID", o.Slug, 1)
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
	return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":pageID", page.Slug, 1)
}

// IsSlugMatch - checks if the slug is a match
func (ds *PageDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/page/")
}
