package datastore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

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

// Iterator - return a list of all Pages, iteration of 1
func (ds *PageCollectionDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	// rule for page 1
	if ctx.Route.PageSize > 0 {

		if !strings.Contains(ctx.Route.URLPath, ":index") {
			panic(fmt.Errorf("Page route is missing an :index. Either add and index placeholder or remove the pageSize\n"))
		}

		fmt.Printf("Pages... page size:%d page total:%d\n", ctx.Route.PageSize, len(ctx.Site.Pages))
		ctx.Route.Pagination = models.Pagination{
			Index: 1,
			Total: (len(ctx.Site.Pages) / ctx.Route.PageSize) + 1,
			Size:  ctx.Route.PageSize,
		}

		// page count
		for pi := 0; pi < ctx.Route.Pagination.Total; pi++ {

			ctx.Route.Pagination.Index = pi + 1
			ctx.Route.Pagination.PreviousURL = ""
			ctx.Route.Pagination.NextURL = ""

			path := strings.Replace(ctx.Route.URLPath, ":index",
				strconv.Itoa(ctx.Route.Pagination.Index), 1)

			if pi > 0 {
				ctx.Route.Pagination.PreviousURL =
					strings.Replace(ctx.Route.URLPath, ":index",
						strconv.Itoa(ctx.Route.Pagination.Index-1), 1)
			}

			if pi < ctx.Route.Pagination.Total-1 {
				ctx.Route.Pagination.NextURL =
					strings.Replace(ctx.Route.URLPath, ":index",
						strconv.Itoa(ctx.Route.Pagination.Index+1), 1)
			}

			startIndex := pi * ctx.Route.PageSize
			endIndex := ((pi * ctx.Route.PageSize) + ctx.Route.PageSize) - 1
			if endIndex >= len(ctx.Site.Pages) {
				endIndex = len(ctx.Site.Pages) - 1
			}

			clonedPages := make([]*models.Page, endIndex-startIndex+1)
			for i := startIndex; i <= endIndex; i++ {
				clonedPages[i-startIndex] = transformPage(ctx.Site.Pages[i])
			}

			vars := make(jet.VarMap)
			vars.Set("pages", clonedPages)
			vars.Set("pagination", ctx.Route.Pagination)
			vars.Set("site", ctx.Site)
			renderer.Render(ctx.Route, ctx.RoutePrefix+path, vars)
		}
	} else {

		ctx.Route.Pagination = models.Pagination{
			Index: 1,
			Total: len(ctx.Site.Pages),
			Size:  len(ctx.Site.Pages),
		}

		clonedPages := make([]*models.Page, len(ctx.Site.Pages))
		for i, f := range ctx.Site.Pages {
			clonedPages[i] = transformPage(f)
		}

		vars := make(jet.VarMap)
		vars.Set("pages", clonedPages)
		vars.Set("pagination", ctx.Route.Pagination)
		vars.Set("site", ctx.Site)
		renderer.Render(ctx.Route, ctx.RoutePrefix+ctx.Route.URLPath, vars)
	}
}

// GetRouteForEntity - get the route
func (ds *PageCollectionDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
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
	f.Content = models.ApplyContentTransforms(f.Content)
	return &f
}
