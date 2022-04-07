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

var tvSeasonArgs = []models.RouteArgument{
	{
		Name:        ":showID",
		Description: "ID of the show",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVSeason); ok {
				return strconv.Itoa(o.ShowInfo.ID)
			}
			return ""
		},
	},
	{
		Name:        ":seasonNumber",
		Description: "Season Number",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVSeason); ok {
				return strconv.Itoa(o.SeasonNumber)
			}
			return ""
		},
	},
	{
		Name:        ":slug",
		Description: "Title slug of the show",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVSeason); ok {
				return o.ShowInfo.TitleSlug
			}
			return ""
		},
	},
}

// TVSeasonDataSource - single tv season datasource
// Supports slugs in the /tv/:tvID/season/:seasonID and /tv/:title_slug
type TVSeasonDataSource struct {
}

// GetName - name of the datasource
func (ds *TVSeasonDataSource) GetName() string {
	return "TVSeason"
}

// GetRouteArguments returns the available route arguments
func (ds *TVSeasonDataSource) GetRouteArguments() []models.RouteArgument {
	return tvSeasonArgs
}

// GetEntityType - Get the entity type
func (ds *TVSeasonDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.TVSeason{})
}

// Iterator - loop over each season
func (ds *TVSeasonDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for i := 0; i < len(ctx.Site.TVSeasons); i++ {
		f := ctx.Site.TVSeasons[i]

		data.Set("tvseason", f)

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
func (ds *TVSeasonDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.URLPath, tvSeasonArgs, entity)
}

// GetPartialRouteForEntity - get the partial route
func (ds *TVSeasonDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.PartialURLPath, tvSeasonArgs, entity)
}

// GetRouteForSlug - get the route
func (ds *TVSeasonDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {

	// supports having tv/:slug/season/:seasonNumber, or any params: :showID, seasonNumber, or :slug
	tvSeason, found := ctx.Site.TVSeasons.FindTVSeasonBySlug(slug)
	if found {
		return ds.GetRouteForEntity(ctx, tvSeason)
	}
	return fmt.Sprintf("ERR(%s)", slug)
}

// IsSlugMatch - checks if the slug is a match
func (ds *TVSeasonDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/tv/") && strings.Contains(slug, "/season/") && !strings.Contains(slug, "/episode/")
}

// IsValid checks for any validation errors
func (ds *TVSeasonDataSource) IsValid(route *models.Route) error {
	return nil
}
