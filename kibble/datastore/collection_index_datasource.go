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
	"reflect"

	"kibble/models"

	"github.com/CloudyKit/jet"
)

// CollectionIndexDataSource - a list of all Collections
type CollectionIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *CollectionIndexDataSource) GetName() string {
	return "CollectionIndex"
}

// GetEntityType - Get the entity type
func (ds *CollectionIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Collection{})
}

// Iterator - return a list of all Collections, iteration of 1
func (ds *CollectionIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	clonedCollections := make([]*models.Collection, len(ctx.Site.Collections))
	for i, f := range ctx.Site.Collections {
		clonedCollections[i] = &f
	}

	vars := make(jet.VarMap)
	vars.Set("collections", clonedCollections)
	vars.Set("site", ctx.Site)
	return renderer.Render(ctx.Route.TemplatePath, ctx.RoutePrefix+ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *CollectionIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *CollectionIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *CollectionIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

// GetRouteArguments returns the available route arguments
func (ds *CollectionIndexDataSource) GetRouteArguments() []models.RouteArgument {
	return indexArgs
}
