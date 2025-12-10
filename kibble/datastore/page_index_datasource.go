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
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
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

type pageIndexDataSourceOptions struct {
	PageTypes []string `json:"pageTypes"`
	SortBy    []string `json:"sortBy"`
}

func (opts pageIndexDataSourceOptions) IsAllowedPageType(pageType string) bool {
	for _, pt := range opts.PageTypes {
		if pt == pageType {
			return true
		}
	}

	return len(opts.PageTypes) == 0
}

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

	var options pageIndexDataSourceOptions
	if len(ctx.Route.Options) > 0 {
		err := json.Unmarshal(ctx.Route.Options, &options)
		if err != nil {
			panic(fmt.Errorf("unable to parse datasource options: %w", err))
		}
	}

	pages := make(models.Pages, 0, len(ctx.Site.Pages))
	for _, page := range ctx.Site.Pages {
		if !options.IsAllowedPageType(page.PageType) {
			continue
		}
		pages = append(pages, page)
	}

	if len(options.SortBy) > 0 {
		sortPages(pages, ParseSortKeys(options.SortBy))
	}

	// rule for page 1
	if ctx.Route.PageSize > 0 {

		if !strings.Contains(ctx.Route.URLPath, ":index") {
			panic(fmt.Errorf("Page route is missing an :index. Either add and index placeholder or remove the pageSize"))
		}

		ctx.Route.Pagination = models.Pagination{
			Index: 1,
			Total: (len(pages) / ctx.Route.PageSize) + 1,
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
			if endIndex >= len(pages) {
				endIndex = len(pages) - 1
			}

			clonedPages := make([]*models.Page, endIndex-startIndex+1)
			for i := startIndex; i <= endIndex; i++ {
				clonedPages[i-startIndex] = transformPage(pages[i])
			}

			vars := make(jet.VarMap)
			vars.Set("pages", clonedPages)
			vars.Set("pagination", ctx.Route.Pagination)
			vars.Set("site", ctx.Site)
			errCount += renderer.Render(ctx.Route.TemplatePath, ctx.RoutePrefix+path, vars)
		}
	} else {

		ctx.Route.Pagination = models.Pagination{
			Index: 1,
			Total: len(pages),
			Size:  len(pages),
		}

		clonedPages := make([]*models.Page, len(pages))
		for i, f := range pages {
			clonedPages[i] = transformPage(f)
		}

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

func transformPage(f models.Page) *models.Page {
	f.Content = models.ApplyContentTransforms(f.Content)

	for i := 0; i < len(f.PageCollections); i++ {
		f.PageCollections[i].Description = models.ApplyContentTransforms(f.PageCollections[i].Description)
	}

	return &f
}

type sortKey struct {
	Field string
	Desc  bool
}

func ParseSortKeys(raw []string) []sortKey {
	keys := make([]sortKey, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		key := sortKey{Field: parts[0]}
		if len(parts) == 2 && strings.EqualFold(parts[1], "desc") {
			key.Desc = true
		}
		keys = append(keys, key)
	}
	return keys
}

func sortPages(pages models.Pages, keys []sortKey) {
	// Apply from last to first so earlier keys win (like SQL ORDER BY)
	for i := len(keys) - 1; i >= 0; i-- {
		key := keys[i]

		sort.SliceStable(pages, func(i, j int) bool {
			pi, pj := pages[i], pages[j]

			switch key.Field {
			case "published_date":
				// handle equal first so we don't randomly flip
				if pi.PublishedDate.Equal(pj.PublishedDate) {
					return false
				}
				if key.Desc {
					return pi.PublishedDate.After(pj.PublishedDate)
				}
				return pi.PublishedDate.Before(pj.PublishedDate)

			case "title":
				if pi.Title == pj.Title {
					return false
				}
				if key.Desc {
					return pi.Title > pj.Title
				}
				return pi.Title < pj.Title

			default:
				// unknown field → no-op for this pass
				return false
			}
		})
	}
}
