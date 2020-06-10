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

	"github.com/CloudyKit/jet"
	"kibble/models"
)

// BundleIndexDataSource - a list of all bundles
type BundleIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *BundleIndexDataSource) GetName() string {
	return "BundleIndex"
}

// GetEntityType - Get the entity type
func (ds *BundleIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Bundle{})
}

// Iterator - return a list of all bundles, iteration of 1
func (ds *BundleIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) int {

	clonedBundles := make([]*models.Bundle, len(ctx.Site.Bundles))
	for i, f := range ctx.Site.Bundles {
		clonedBundles[i] = transformBundle(f)
	}

	vars := make(jet.VarMap)
	vars.Set("bundles", clonedBundles)
	vars.Set("site", ctx.Site)
	return renderer.Render(ctx.Route.TemplatePath, ctx.RoutePrefix+ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *BundleIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *BundleIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *BundleIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformBundle(f models.Bundle) *models.Bundle {
	f.Description = models.ApplyContentTransforms(f.Description)
	return &f
}

// GetRouteArguments returns the available route arguments
func (ds *BundleIndexDataSource) GetRouteArguments() []models.RouteArgument {
	return indexArgs
}
