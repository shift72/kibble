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

// TVShowIndexDataSource - a list of all films
type TVShowIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *TVShowIndexDataSource) GetName() string {
	return "TVShowIndex"
}

// GetEntityType - Get the entity type
func (ds *TVShowIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.TVShow{})
}

// Iterator - return a list of all tv seasons, iteration of 1
func (ds *TVShowIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	vars := make(jet.VarMap)
	vars.Set("tvshows", ctx.Site.TVShows)
	vars.Set("site", ctx.Site)
	return renderer.Render(ctx.Route.TemplatePath, ctx.RoutePrefix+ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *TVShowIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *TVShowIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *TVShowIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformTVShow(f *models.TVShow) *models.TVShow {
	f.Overview = models.ApplyContentTransforms(f.Overview)
	return f
}

// GetRouteArguments returns the available route arguments
func (ds *TVShowIndexDataSource) GetRouteArguments() []models.RouteArgument {
	return indexArgs
}
