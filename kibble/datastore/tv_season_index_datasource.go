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
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// TVSeasonIndexDataSource - a list of all films
type TVSeasonIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *TVSeasonIndexDataSource) GetName() string {
	return "TVSeasonIndex"
}

// GetEntityType - Get the entity type
func (ds *TVSeasonIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.TVSeason{})
}

// Iterator - return a list of all tv seasons, iteration of 1
func (ds *TVSeasonIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	cloned := make([]*models.TVSeason, len(ctx.Site.TVSeasons))
	for i, s := range ctx.Site.TVSeasons {
		cloned[i] = transformTVSeason(s)
	}

	vars := make(jet.VarMap)
	vars.Set("tvseasons", cloned)
	vars.Set("site", ctx.Site)
	return renderer.Render(ctx.Route.TemplatePath, ctx.RoutePrefix+ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *TVSeasonIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *TVSeasonIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *TVSeasonIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

func transformTVSeason(f models.TVSeason) *models.TVSeason {
	f.Overview = models.ApplyContentTransforms(f.Overview)

	if f.ShowInfo != nil {
		f.ShowInfo.Overview = models.ApplyContentTransforms(f.ShowInfo.Overview)
	}

	for i, e := range f.Episodes {
		ov := models.ApplyContentTransforms(e.Overview)
		// ranges create a copy of the array, so we need to set the original object
		f.Episodes[i].Overview = ov
	}

	return &f
}

// GetRouteArguments returns the available route arguments
func (ds *TVSeasonIndexDataSource) GetRouteArguments() []models.RouteArgument {
	return indexArgs
}
