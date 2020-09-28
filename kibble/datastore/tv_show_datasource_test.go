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

	"github.com/stretchr/testify/assert"

	"kibble/models"
	"kibble/test"
)

func createTestTVShow() (models.RenderContext, *models.Route) {
	return createTestTVShowWithCustomURLPath("/tv/:slug")
}

func createTestTVShowWithCustomURLPath(urlPath string) (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      urlPath,
		TemplatePath: "tv/item.jet",
		DataSource:   "TVShow",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			TVShows: models.TVShowCollection{
				&models.TVShow{
					ID:        3,
					Title:     "Breaking Bad",
					TitleSlug: "breaking-bad",
					Slug:      "/tv/3",
				},
			},
		},
	}

	return ctx, r
}

func TestTVShowGetRouteForSlug(t *testing.T) {
	var TVShowDS TVShowDataSource

	ctx, _ := createTestTVShow()

	route := TVShowDS.GetRouteForSlug(ctx, "/tv/3")

	assert.Equal(t, "/fr/tv/breaking-bad", route)
}

func TestTVShowIsSlugMatch(t *testing.T) {
	var TVShowDS TVShowDataSource

	assert.True(t, TVShowDS.IsSlugMatch("/tv/123"))
	assert.False(t, TVShowDS.IsSlugMatch("/tv/123/season/2"))
}

func TestTVShowGetRouteForMissingSlug(t *testing.T) {
	var TVShowDS TVShowDataSource

	ctx, _ := createTestTVShow()

	route := TVShowDS.GetRouteForSlug(ctx, "/tv/999")

	assert.Equal(t, "ERR(/tv/999)", route)
}

func TestTVShowGetRouteForInvalidSlug(t *testing.T) {
	var TVShowDS TVShowDataSource

	ctx, _ := createTestTVShow()

	route := TVShowDS.GetRouteForSlug(ctx, "/tv/a")

	assert.Equal(t, "ERR(/tv/a)", route)
}

func TestTVShowGetRouteWithIDForSlug(t *testing.T) {
	var TVShowDS TVShowDataSource

	ctx, _ := createTestTVShowWithCustomURLPath("/tv/:showID.html")

	route := TVShowDS.GetRouteForSlug(ctx, "/tv/3")

	assert.Equal(t, "/fr/tv/3.html", route)
}

func TestRenderTVShow(t *testing.T) {
	var ds TVShowDataSource

	renderer := &test.MockRenderer{}

	ctx, _ := createTestTVShow()

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/fr/tv/breaking-bad", renderer.FilePath)
	assert.Equal(t, "tv/item.jet", renderer.TemplatePath)
}

func TestPartialRenderTVShow(t *testing.T) {
	var ds TVShowDataSource

	renderer := &test.MockRenderer{}

	ctx, r := createTestTVShow()
	r.PartialTemplatePath = "/tv/partial.jet"
	r.PartialURLPath = "/partials/tv/:showID.html"

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/fr/partials/tv/3.html", renderer.FilePath)
	assert.Equal(t, "/tv/partial.jet", renderer.TemplatePath)
}
