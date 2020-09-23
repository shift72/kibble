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
	"testing"

	"kibble/models"
	"kibble/test"

	"github.com/nicksnyder/go-i18n/i18n"
)

func TestRenderingShowSlug(t *testing.T) {

	Init()

	r := &models.Route{
		Name:         "tvShowItem",
		URLPath:      "/tv-show/:slug",
		TemplatePath: "tvshow/item.jet",
		DataSource:   "TVShow",
	}

	cfg := models.Config{
		Routes: []models.Route{*r},
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			TVShows: models.TVShowCollection{
				&models.TVShow{
					ID:        123,
					Slug:      "/tv/123",
					TitleSlug: "the-big-lebowski",
				},
			},
		},
	}

	routeRegistry := models.NewRouteRegistryFromConfig(&cfg)

	view := models.CreateTemplateView(routeRegistry, i18n.IdentityTfunc(), &ctx, "./templates")

	tem, _ := view.LoadTemplate("", "{{ routeToSlug(site.TVShows[0].Slug) }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	var ds TVShowDataSource
	ds.Iterator(ctx, renderer)

	if renderer.Result.Output() != "/fr/tv-show/the-big-lebowski" {
		t.Errorf("Unexpected output. `%s`", renderer.Result.Output())
	}
}
