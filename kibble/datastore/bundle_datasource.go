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

// Iterator - loop over each Bundle
func (ds *BundleDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for _, p := range ctx.Site.Bundles {
		c := transformBundle(p)
		data.Set("bundle", c)

		// normal bundle pages
		filePath := ds.GetRouteForEntity(ctx, &p)
		errCount += renderer.Render(ctx.Route, filePath, data)

		// bundle partials
		if ctx.Route.HasPartial() {
			partialFilePath := ds.GetPartialRouteForEntity(ctx, &p)
			errCount += renderer.Render(ctx.Route, partialFilePath, data)
		}
	}
	return
}

// GetRouteForEntity - get the route
func (ds *BundleDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Bundle)
	if ok {
		s := strings.Replace(ctx.Route.URLPath, ":slug", o.TitleSlug, 1)
		s = strings.Replace(s, ":bundleID", strconv.Itoa(o.ID), 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
}

// GetPartialRouteForEntity - get the partial route
func (ds *BundleDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Bundle)
	if ok {
		s := strings.Replace(ctx.Route.PartialURLPath, ":slug", o.TitleSlug, 1)
		s = strings.Replace(s, ":bundleID", strconv.Itoa(o.ID), 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *BundleDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	bundleID, ok := utils.ParseIntFromSlug(slug, 2)
	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}
	bundle, err := ctx.Site.Bundles.FindBundleByID(bundleID)

	if err != nil {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	return ds.GetRouteForEntity(ctx, bundle)
}

// IsSlugMatch - checks if the slug is a match
func (ds *BundleDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/bundle/")
}
