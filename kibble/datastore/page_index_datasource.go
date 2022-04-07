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

var indexArgs = []models.RouteArgument{
	{
		Name:        ":index",
		Description: "Index of the page",
		GetValue: func(entity interface{}) string {
			return ""
		},
	},
}

// PageIndexDataSource - a list of all Pages
type PageIndexDataSource struct{}

// GetName - returns the name of the datasource
func (ds *PageIndexDataSource) GetName() string {
	return "PageIndex"
}

// GetRouteArguments returns the available route arguments
func (ds *PageIndexDataSource) GetRouteArguments() []models.RouteArgument {
	return indexArgs
}

// GetEntityType - Get the entity type
func (ds *PageIndexDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf([]models.Page{})
}

// Iterator - return a list of all Pages, iteration of 1
func (ds *PageIndexDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {
	// rule for page 1
	if ctx.Route.PageSize > 0 {

		if !strings.Contains(ctx.Route.URLPath, ":index") {
			panic(fmt.Errorf("Page route is missing an :index. Either add and index placeholder or remove the pageSize"))
		}

		ctx.Route.Pagination = models.Pagination{
			Index: 1,
			Total: (len(ctx.Site.Pages) / ctx.Route.PageSize) + 1,
			Size:  ctx.Route.PageSize,
		}

		// page count
		for pi := 0; pi < ctx.Route.Pagination.Total; pi++ {

			ctx.Route.Pagination.Index = pi + 1
			ctx.Route.Pagination.PreviousURL = ""
			ctx.Route.Pagination.NextURL = ""

			path := strings.Replace(ctx.Route.URLPath, ":index",
				strconv.Itoa(ctx.Route.Pagination.Index), 1)

			if pi > 0 {
				ctx.Route.Pagination.PreviousURL =
					strings.Replace(ctx.Route.URLPath, ":index",
						strconv.Itoa(ctx.Route.Pagination.Index-1), 1)
			}

			if pi < ctx.Route.Pagination.Total-1 {
				ctx.Route.Pagination.NextURL =
					strings.Replace(ctx.Route.URLPath, ":index",
						strconv.Itoa(ctx.Route.Pagination.Index+1), 1)
			}

			startIndex := pi * ctx.Route.PageSize
			endIndex := ((pi * ctx.Route.PageSize) + ctx.Route.PageSize) - 1
			if endIndex >= len(ctx.Site.Pages) {
				endIndex = len(ctx.Site.Pages) - 1
			}

			clonedPages := make([]*models.Page, endIndex-startIndex+1)
			vars := make(jet.VarMap)
			vars.Set("pages", clonedPages)
			vars.Set("pagination", ctx.Route.Pagination)
			vars.Set("site", ctx.Site)
			errCount += renderer.Render(ctx.Route.TemplatePath, ctx.RoutePrefix+path, vars)
		}
	} else {

		ctx.Route.Pagination = models.Pagination{
			Index: 1,
			Total: len(ctx.Site.Pages),
			Size:  len(ctx.Site.Pages),
		}

		clonedPages := make([]*models.Page, len(ctx.Site.Pages))
		vars := make(jet.VarMap)
		vars.Set("pages", clonedPages)
		vars.Set("pagination", ctx.Route.Pagination)
		vars.Set("site", ctx.Site)
		errCount += renderer.Render(ctx.Route.TemplatePath, ctx.RoutePrefix+ctx.Route.URLPath, vars)
	}

	return
}

// GetRouteForEntity - get the route
func (ds *PageIndexDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return ctx.RoutePrefix + ctx.Route.URLPath
}

// GetRouteForSlug - get the route
func (ds *PageIndexDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return models.ErrDataSource
}

// IsSlugMatch - is the slug a match
func (ds *PageIndexDataSource) IsSlugMatch(slug string) bool {
	return false
}
