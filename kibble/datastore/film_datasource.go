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

var filmArgs = []models.RouteArgument{
	{
		Name:        ":filmID",
		Description: "ID of the collection",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.Film); ok {
				return strconv.Itoa(o.ID)
			}
			return ""
		},
	},
	{
		Name:        ":slug",
		Description: "Slug of the collection",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.Film); ok {
				return o.TitleSlug
			}
			return ""
		},
	},
}

// FilmDataSource - single film datasource
// Supports slugs in the /film/:filmID and /film/:title_slug
type FilmDataSource struct {
}

// GetName - name of the datasource
func (ds *FilmDataSource) GetName() string {
	return "Film"
}

// GetRouteArguments returns the available route arguments
func (ds *FilmDataSource) GetRouteArguments() []models.RouteArgument {
	return filmArgs
}

// GetEntityType - Get the entity type
func (ds *FilmDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.Film{})
}

// Iterator - loop over each film
func (ds *FilmDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {
	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for _, f := range ctx.Site.Films {
		data.Set("film", f)

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
func (ds *FilmDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.URLPath, filmArgs, entity)
}

// GetPartialRouteForEntity - get the partial route
func (ds *FilmDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.PartialURLPath, filmArgs, entity)
}

// GetRouteForSlug - get the route
func (ds *FilmDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	film, ok := ctx.Site.Films.FindFilmBySlug(slug)
	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	return ds.GetRouteForEntity(ctx, film)
}

// IsSlugMatch - checks if the slug is a match
func (ds *FilmDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/film/")
}

// IsValid checks for any validation errors
func (ds *FilmDataSource) IsValid(route *models.Route) error {
	return nil
}
