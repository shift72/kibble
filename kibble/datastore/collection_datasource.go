//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package datastore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"kibble/models"
	"kibble/utils"
)

var collectionArgs = []models.RouteArgument{
	models.RouteArgument{
		Name:        ":collectionID",
		Description: "ID of the collection",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.Collection); ok {
				return strconv.Itoa(o.ID)
			}
			return models.ErrDataSource
		},
	},
	models.RouteArgument{
		Name:        ":slug",
		Description: "Slug of the collection",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.Collection); ok {
				return o.TitleSlug
			}
			return models.ErrDataSource
		},
	},
}

// CollectionDataSource - single Collection datasource
type CollectionDataSource struct {
}

// GetName - name of the datasource
func (ds *CollectionDataSource) GetName() string {
	return "Collection"
}

// GetRouteArguments returns the available route arguments
func (ds *CollectionDataSource) GetRouteArguments() []models.RouteArgument {
	return collectionArgs
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
		errCount += renderer.Render(ctx.Route.TemplatePath, filePath, data)

		if ctx.Route.HasPartial() {
			partialFilePath := ds.GetPartialRouteForEntity(ctx, &p)
			errCount += renderer.Render(ctx.Route.PartialTemplatePath, partialFilePath, data)
		}
	}

	return
}

// GetRouteForEntity - get the route
func (ds *CollectionDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.URLPath, collectionArgs, entity)
}

// GetPartialRouteForEntity - get the partial route
func (ds *CollectionDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.PartialURLPath, collectionArgs, entity)
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

// IsValid checks for any validation errors
func (ds *CollectionDataSource) IsValid(route *models.Route) error {
	return nil
}
