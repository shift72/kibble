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
	"kibble/utils"

	"github.com/CloudyKit/jet"
)

var bundleArgs = []models.RouteArgument{
	{
		Name:        ":bundleID",
		Description: "ID of the bundle",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.Bundle); ok {
				return strconv.Itoa(o.ID)
			}
			return ""
		},
	},
	{
		Name:        ":slug",
		Description: "Slug of the bundle",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.Bundle); ok {
				return o.TitleSlug
			}
			return ""
		},
	},
}

// BundleDataSource - single Bundle datasource
type BundleDataSource struct {
}

// GetName - name of the datasource
func (ds *BundleDataSource) GetName() string {
	return "Bundle"
}

// GetRouteArguments returns the available route arguments
func (ds *BundleDataSource) GetRouteArguments() []models.RouteArgument {
	return bundleArgs
}

// GetEntityType - Get the entity type
func (ds *BundleDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.Bundle{})
}

// Iterator - loop over each Bundle
func (ds *BundleDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for _, p := range ctx.Site.Bundles {
		c := transformBundle(p)
		data.Set("bundle", c)

		// normal bundle pages
		filePath := ds.GetRouteForEntity(ctx, &p)
		errCount += renderer.Render(ctx.Route.TemplatePath, filePath, data)

		// bundle partials
		if ctx.Route.HasPartial() {
			partialFilePath := ds.GetPartialRouteForEntity(ctx, &p)
			errCount += renderer.Render(ctx.Route.PartialTemplatePath, partialFilePath, data)
		}
	}
	return
}

// GetRouteForEntity - get the route
func (ds *BundleDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.URLPath, bundleArgs, entity)
}

// GetPartialRouteForEntity - get the partial route
func (ds *BundleDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.PartialURLPath, bundleArgs, entity)
}

// GetRouteForSlug - get the route
func (ds *BundleDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	bundleID, ok := utils.ParseIntFromSlug(slug, 2)
	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}
	bundle, err := ctx.Site.Bundles.FindBundleByID(bundleID)

	if err != nil {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	return ds.GetRouteForEntity(ctx, bundle)
}

// IsSlugMatch - checks if the slug is a match
func (ds *BundleDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/bundle/")
}
