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
	"github.com/indiereign/shift72-kibble/kibble/utils"
)

var pageArgs = []models.RouteArgument{
	models.RouteArgument{
		Name:        ":pageID",
		Description: "ID of the page",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.Page); ok {
				return strconv.Itoa(o.ID)
			}
			return ""
		},
	},
	models.RouteArgument{
		Name:        ":slug",
		Description: "Slug of the page",
		GetValue: func(entity interface{}) string {
			if o, ok := entity.(*models.Page); ok {
				return o.TitleSlug
			}
			return ""
		},
	},
}

// PageDataSource - single Page datasource
type PageDataSource struct {
}

// GetName - name of the datasource
func (ds *PageDataSource) GetName() string {
	return "Page"
}

// GetRouteArguments returns the available route arguments
func (ds *PageDataSource) GetRouteArguments() []models.RouteArgument {
	return pageArgs
}

// GetEntityType - Get the entity type
func (ds *PageDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.Page{})
}

// Iterator - loop over each Page
func (ds *PageDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	lang := ""
	if !ctx.Language.IsDefault {
		lang = ctx.Language.Code
	}

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)
	data.Set("currentLanguage", lang)

	for _, p := range ctx.Site.Pages {

		availableI18n := make([]string, 0)
		availableI18n = append(availableI18n, ctx.Site.SiteConfig.DefaultLanguage)

		// Convert (e.g) '/fr/page/about-us' to '/page/about-us' so we can keep track of the original URLs
		defaultURLPath := strings.Replace(ds.GetRouteForEntity(ctx, &p), "/"+lang+"/", "/", 1)
		data.Set("defaultUrlPath", defaultURLPath)

		// Check there's a translation for this language and use that if possible
		if len(p.Translations) > 0 {
			for _, translation := range p.Translations {
				availableI18n = append(availableI18n, translation.Language)
				if lang != "" && translation.Language == lang {
					p = translation.Page

				}
			}
		}

		p.AvailableI18n = availableI18n
		data.Set("page", transformPage(p))

		// render page endpoints first
		// don't render external pages
		if p.PageType != "external" {
			templatePath := strings.Replace(ctx.Route.TemplatePath, ":type", p.PageType, 1)
			errCount += renderer.Render(templatePath, ds.GetRouteForEntity(ctx, &p), data)
		}

		// now partial end points
		if ctx.Route.HasPartial() {
			templatePath := strings.Replace(ctx.Route.PartialTemplatePath, ":type", p.PageType, 1)
			errCount += renderer.Render(templatePath, ds.GetPartialRouteForEntity(ctx, &p), data)
		}
	}

	return
}

// GetRouteForEntity - get the route
func (ds *PageDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Page)
	if ok {

		switch o.PageType {
		// special case for the homepage
		case "homepage":
			return ctx.RoutePrefix + "/"
		case "external":
			// special case: map an old static url to a new one
			if strings.Contains(o.URL, "/#!/") &&
				strings.HasPrefix(o.URL, ctx.Site.SiteConfig.SiteURL) {
				i := strings.Index(o.URL, "/#!/") + 3
				return o.URL[i:len(o.URL)]
			}
			return o.URL
		default:
			return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.URLPath, pageArgs, entity)
		}
	}
	return models.ErrDataSource
}

// GetPartialRouteForEntity - get the partial route
func (ds *PageDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ReplaceURLArgumentsWithEntityValues(ctx.RoutePrefix, ctx.Route.PartialURLPath, pageArgs, entity)
}

// GetRouteForSlug - get the route
func (ds *PageDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {

	pageID, ok := utils.ParseIntFromSlug(slug, 2)
	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	page, ok := ctx.Site.Pages.FindPageByID(pageID)

	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	return ds.GetRouteForEntity(ctx, page)
}

// IsSlugMatch - checks if the slug is a match
func (ds *PageDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/page/")
}

// IsValid checks for any validation errors
func (ds *PageDataSource) IsValid(route *models.Route) error {
	return nil
}
