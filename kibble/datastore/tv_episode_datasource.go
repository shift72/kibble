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

var tvEpisodeArgs = []models.RouteArgument{
	models.RouteArgument{
		Name:        ":showID",
		Description: "ID of the show",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVEpisode); ok {
				return strconv.Itoa(o.Season.ShowInfo.ID)
			}
			return ""
		},
	},
	models.RouteArgument{
		Name:        ":seasonNumber",
		Description: "Season Number",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVEpisode); ok {
				return strconv.Itoa(o.Season.SeasonNumber)
			}
			return ""
		},
	},
	models.RouteArgument{
		Name:        ":slug",
		Description: "Title slug of the show",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVEpisode); ok {
				return o.Season.ShowInfo.TitleSlug
			}
			return ""
		},
	},
	models.RouteArgument{
		Name:        ":episodeNumber",
		Description: "Episode Number",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVEpisode); ok {
				return strconv.Itoa(o.EpisodeNumber)
			}
			return ""
		},
	},
	models.RouteArgument{
		Name:        ":episodeSlug",
		Description: "Title slug of the episode",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.TVEpisode); ok {
				return o.TitleSlug
			}
			return ""
		},
	},
}

// TVEpisodeDataSource - single tv episode datasource
// Supports slugs in the /tv/:tvID/season/:seasonNumber/episode/:episodeNumber format
type TVEpisodeDataSource struct {
}

// GetName - name of the datasource
func (ds *TVEpisodeDataSource) GetName() string {
	return "TVEpisode"
}

// GetRouteArguments returns the available route arguments
func (ds *TVEpisodeDataSource) GetRouteArguments() []models.RouteArgument {
	return tvEpisodeArgs
}

// GetEntityType - Get the entity type
func (ds *TVEpisodeDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.TVEpisode{})
}

// Iterator - loop over each film
func (ds *TVEpisodeDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for i := 0; i < len(ctx.Site.TVEpisodes); i++ {
		episode := ctx.Site.TVEpisodes[i]

		data.Set("tvepisode", transformEpisode(episode))
		if len(ctx.Route.TemplatePath) > 0 {
			filePath := ds.GetRouteForEntity(ctx, episode)
			errCount += renderer.Render(ctx.Route.TemplatePath, filePath, data)
		}

		if ctx.Route.HasPartial() {
			partialFilePath := ds.GetPartialRouteForEntity(ctx, episode)
			errCount += renderer.Render(ctx.Route.PartialTemplatePath, partialFilePath, data)
		}
	}

	return
}

// GetRouteForEntity - get the route
func (ds *TVEpisodeDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.URLPath, tvEpisodeArgs, entity)
}

// GetPartialRouteForEntity - get the partial route
func (ds *TVEpisodeDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.PartialURLPath, tvEpisodeArgs, entity)
}

// GetRouteForSlug - get the route
func (ds *TVEpisodeDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {

	// supports having tv/:slug/season/:seasonNumber, or any params: :showID, seasonNumber, or :slug
	episode, found := ctx.Site.TVEpisodes.FindTVEpisodeBySlug(slug)
	if found {
		return ds.GetRouteForEntity(ctx, episode)
	}
	return fmt.Sprintf("ERR(%s)", slug)
}

// IsSlugMatch - checks if the slug is a match
func (ds *TVEpisodeDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/tv/") && strings.Contains(slug, "/season/") && strings.Contains(slug, "/episode/")
}

// IsValid checks for any validation errors
func (ds *TVEpisodeDataSource) IsValid(route *models.Route) error {
	return nil
}

// transformEpisode applies any content transforms to the episodes overview field
func transformEpisode(e *models.TVEpisode) *models.TVEpisode {
	ov := models.ApplyContentTransforms(e.Overview)
	// ranges create a copy of the array, so we need to set the original object
	e.Overview = ov

	return e
}
