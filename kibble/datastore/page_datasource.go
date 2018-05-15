package datastore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/utils"
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
func (ds *PageDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for _, p := range ctx.Site.Pages {

		data.Set("page", transformPage(p))

		// render page endpoints first
		// don't render external pages
		if p.PageType != "external" {
			templatePath := strings.Replace(ctx.Route.TemplatePath, ":type", p.PageType, 1)
			errCount += renderer.Render(templatePath, ds.GetRouteForEntity(ctx, &p), data)
		}

		// now partial end points
		if ctx.Route.HasPartial() {
			templatePath := strings.Replace(ctx.Route.PartialTemplatePath, ":type", p.PageType, 1)
			errCount += renderer.Render(templatePath, ds.GetPartialRouteForEntity(ctx, &p), data)
		}

	}

	return
}

// GetRouteForEntity - get the route
func (ds *PageDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Page)
	if ok {

		switch o.PageType {
		// special case for the homepage
		case "homepage":
			return ctx.RoutePrefix + "/"
		case "external":
			// special case: map an old static url to a new one
			if strings.Contains(o.URL, "/#!/") &&
				strings.HasPrefix(o.URL, ctx.Site.SiteConfig.SiteURL) {
				i := strings.Index(o.URL, "/#!/") + 3
				return o.URL[i:len(o.URL)]
			}
			return o.URL
		default:
			s := strings.Replace(ctx.Route.URLPath, ":slug", o.TitleSlug, 1)
			s = strings.Replace(s, ":pageID", strconv.Itoa(o.ID), 1)
			return ctx.RoutePrefix + s
		}
	}
	return models.ErrDataSource
}

// GetPartialRouteForEntity - get the partial route
func (ds *PageDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Page)
	if ok {
		s := strings.Replace(ctx.Route.PartialURLPath, ":slug", o.TitleSlug, 1)
		s = strings.Replace(s, ":pageID", strconv.Itoa(o.ID), 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *PageDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {

	pageID, ok := utils.ParseIntFromSlug(slug, 2)
	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	page, ok := ctx.Site.Pages.FindPageByID(pageID)

	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	return ds.GetRouteForEntity(ctx, page)
}

// IsSlugMatch - checks if the slug is a match
func (ds *PageDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/page/")
}
