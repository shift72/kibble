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
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// TVSeasonDataSource - single tv season datasource
// Supports slugs in the /tv/:tvID/season/:seasonID and /tv/:title_slug
type TVSeasonDataSource struct{}

// GetName - name of the datasource
func (ds *TVSeasonDataSource) GetName() string {
	return "TVSeason"
}

// GetEntityType - Get the entity type
func (ds *TVSeasonDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.TVSeason{})
}

// Iterator - loop over each film
func (ds *TVSeasonDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for _, f := range ctx.Site.TVSeasons {
		data.Set("tvseason", transformTVSeason(f))

		filePath := ds.GetRouteForEntity(ctx, &f)
		errCount += renderer.Render(ctx.Route.TemplatePath, filePath, data)

		if ctx.Route.HasPartial() {
			partialFilePath := ds.GetPartialRouteForEntity(ctx, &f)
			errCount += renderer.Render(ctx.Route.PartialTemplatePath, partialFilePath, data)
		}
	}

	return
}

// GetRouteForEntity - get the route
func (ds *TVSeasonDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.TVSeason)
	if ok {
		s := strings.Replace(ctx.Route.URLPath, ":slug", o.ShowInfo.TitleSlug, 1)
		s = strings.Replace(s, ":seasonNumber", strconv.Itoa(o.SeasonNumber), 1)
		s = strings.Replace(s, ":showID", strconv.Itoa(o.ShowInfo.ID), 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
}

// GetPartialRouteForEntity - get the partial route
func (ds *TVSeasonDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.TVSeason)
	if ok {
		s := strings.Replace(ctx.Route.PartialURLPath, ":slug", o.ShowInfo.TitleSlug, 1)
		s = strings.Replace(s, ":seasonNumber", strconv.Itoa(o.SeasonNumber), 1)
		s = strings.Replace(s, ":showID", strconv.Itoa(o.ShowInfo.ID), 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
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
	return strings.HasPrefix(slug, "/tv/") && strings.Contains(slug, "/season/")
}
