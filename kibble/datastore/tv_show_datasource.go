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

	"kibble/models"

	"github.com/CloudyKit/jet"
)

var tvShowArgs = []models.RouteArgument{
	models.RouteArgument{
		Name:        ":showID",
		Description: "ID of the show",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVShow); ok {
				return strconv.Itoa(o.ID)
			}
			return ""
		},
	},
	models.RouteArgument{
		Name:        ":slug",
		Description: "Title slug of the show",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVShow); ok {
				return o.TitleSlug
			}
			return ""
		},
	},
}

// TVShowDataSource - single tv season datasource
// Supports slugs in the /tv/:tvID/season/:seasonID and /tv/:title_slug
type TVShowDataSource struct {
}

// GetName - name of the datasource
func (ds *TVShowDataSource) GetName() string {
	return "TVShow"
}

// GetRouteArguments returns the available route arguments
func (ds *TVShowDataSource) GetRouteArguments() []models.RouteArgument {
	return tvShowArgs
}

// GetEntityType - Get the entity type
func (ds *TVShowDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.TVShow{})
}

// Iterator - loop over each film
func (ds *TVShowDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for i := 0; i < len(ctx.Site.TVShows); i++ {
		f := ctx.Site.TVShows[i]

		data.Set("tvshow", transformTVShow(f))

		filePath := ds.GetRouteForEntity(ctx, f)
		errCount += renderer.Render(ctx.Route.TemplatePath, filePath, data)

		if ctx.Route.HasPartial() {
			partialFilePath := ds.GetPartialRouteForEntity(ctx, f)
			errCount += renderer.Render(ctx.Route.PartialTemplatePath, partialFilePath, data)
		}
	}

	return
}

// GetRouteForEntity - get the route
func (ds *TVShowDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.URLPath, tvShowArgs, entity)
}

// GetPartialRouteForEntity - get the partial route
func (ds *TVShowDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.PartialURLPath, tvShowArgs, entity)
}

// GetRouteForSlug - get the route
func (ds *TVShowDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {

	// supports any params: :slug, :showID
	tvShow, found := ctx.Site.TVShows.FindTVShowBySlug(slug)
	if found {
		return ds.GetRouteForEntity(ctx, tvShow)
	}
	return fmt.Sprintf("ERR(%s)", slug)
}

// IsSlugMatch - checks if the slug is a match
func (ds *TVShowDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/tv/") && !strings.Contains(slug, "/season/")
}

// IsValid checks for any validation errors
func (ds *TVShowDataSource) IsValid(route *models.Route) error {
	return nil
}
