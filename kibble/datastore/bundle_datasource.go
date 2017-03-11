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

// BundleDataSource - single Bundle datasource
type BundleDataSource struct{}

// GetName - name of the datasource
func (ds *BundleDataSource) GetName() string {
	return "Bundle"
}

// GetEntityType - Get the entity type
func (ds *BundleDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.Bundle{})
}

// Query - return a single Bundle
func (ds *BundleDataSource) Query(ctx models.RenderContext, req *http.Request) (jet.VarMap, error) {

	bundleSlug := chi.URLParam(req, "slug")

	p, err := ctx.Site.Bundles.FindBundleBySlug(bundleSlug)
	if err != nil || p == nil {
		return nil, err
	}
	c := transformBundle(*p)

	vars := make(jet.VarMap)
	vars.Set("bundle", c)
	vars.Set("site", ctx.Site)
	return vars, nil
}

// Iterator - loop over each Bundle
func (ds *BundleDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	data := make(jet.VarMap)

	for _, p := range ctx.Site.Bundles {

		filePath := ds.GetRouteForEntity(ctx, &p)

		c := transformBundle(p)

		data.Set("bundle", c)
		data.Set("site", ctx.Site)
		renderer.Render(ctx.Route, filePath, data)
	}
}

// GetRouteForEntity - get the route
func (ds *BundleDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Bundle)
	if ok {
		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":slug", o.TitleSlug, 1)
	}
	return models.DataSourceError
}

// GetRouteForSlug - get the route
func (ds *BundleDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	p := strings.Split(slug, "/")
	bundleID, _ := strconv.Atoi(p[2])
	bundle, _ := ctx.Site.Bundles.FindBundleByID(bundleID)
	return ds.GetRouteForEntity(ctx, bundle)
}

// IsSlugMatch - checks if the slug is a match
func (ds *BundleDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/bundle/")
}
