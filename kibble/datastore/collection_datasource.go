package datastore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
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

	for _, p := range ctx.Site.Collections {

		filePath := ds.GetRouteForEntity(ctx, &p)

		c := transformCollection(p)

		data.Set("collection", c)
		data.Set("site", ctx.Site)
		errCount += renderer.Render(ctx.Route, filePath, data)
	}

	return
}

// GetRouteForEntity - get the route
func (ds *CollectionDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {

	o, ok := entity.(*models.Collection)
	if ok {
		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":slug", o.TitleSlug, 1)
	}
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *CollectionDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	p := strings.Split(slug, "/")
	collectionID, err := strconv.Atoi(p[2])
	if err != nil {
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
