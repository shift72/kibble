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

// CollectionDataSource - single Collection datasource
type CollectionDataSource struct{}

// GetName - name of the datasource
func (ds *CollectionDataSource) GetName() string {
	return "Collection"
}

// GetEntityType - Get the entity type
func (ds *CollectionDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.Collection{})
}

// Iterator - loop over each Collection
func (ds *CollectionDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for _, p := range ctx.Site.Collections {
		c := transformCollection(p)
		data.Set("collection", c)

		filePath := ds.GetRouteForEntity(ctx, &p)
		errCount += renderer.Render(ctx.Route, filePath, data)

		if ctx.Route.HasPartial() {
			route := ctx.Route.Clone()
			route.TemplatePath = ctx.Route.PartialTemplatePath
			partialFilePath := ds.GetPartialRouteForEntity(ctx, &p)
			errCount += renderer.Render(route, partialFilePath, data)
		}
	}

	return
}

// GetRouteForEntity - get the route
func (ds *CollectionDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {

	o, ok := entity.(*models.Collection)
	if ok {
		url := ctx.Route.URLPath
		if strings.Contains(url, ":collectionID") {
			url = strings.Replace(url, ":collectionID", strconv.Itoa(o.ID), 1)
		}

		if strings.Contains(url, ":slug") {
			url = strings.Replace(url, ":slug", o.TitleSlug, 1)
		}

		return ctx.RoutePrefix + url
	}
	return models.ErrDataSource
}

// GetPartialRouteForEntity - get the partial route
func (ds *CollectionDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {

	o, ok := entity.(*models.Collection)
	if ok {
		url := strings.Replace(ctx.Route.PartialURLPath, ":slug", o.TitleSlug, 1)
		url = strings.Replace(url, ":collectionID", strconv.Itoa(o.ID), 1)
		return ctx.RoutePrefix + url
	}
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *CollectionDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	collectionID, ok := utils.ParseIntFromSlug(slug, 2)
	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}
	collection, err := ctx.Site.Collections.FindCollectionByID(collectionID)
	if err != nil {
		return fmt.Sprintf("ERR(%s)", slug)
	}
	return ds.GetRouteForEntity(ctx, collection)
}

// IsSlugMatch - checks if the slug is a match
func (ds *CollectionDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/feature/") ||
		strings.HasPrefix(slug, "/collection/")
}
