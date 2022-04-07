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
	"kibble/models"
	"reflect"

	"github.com/CloudyKit/jet"

	logging "github.com/op/go-logging"
)

// FilmIndexDataSource - a list of all films
type FilmIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *FilmIndexDataSource) GetName() string {
	return "FilmIndex"
}

// GetEntityType - Get the entity type
func (ds *FilmIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Film{})
}

// Iterator - return a list of all films, iteration of 1
func (ds *FilmIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {
	clonedFilms := make([]*models.Film, len(ctx.Site.Films))

	i := 0
	for _, f := range ctx.Site.Films {
		clonedFilms[i] = f
		i++
	}

	vars := make(jet.VarMap)
	vars.Set("films", clonedFilms)
	vars.Set("site", ctx.Site)
	return renderer.Render(ctx.Route.TemplatePath, ctx.RoutePrefix+ctx.Route.URLPath, vars)
}

// GetRouteForEntity - get the route
func (ds *FilmIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *FilmIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *FilmIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}

var log = logging.MustGetLogger("datastore")

// GetRouteArguments returns the available route arguments
func (ds *FilmIndexDataSource) GetRouteArguments() []models.RouteArgument {
	return indexArgs
}
