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

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/test"
)

func createTestTVSeason() (models.RenderContext, *models.Route) {
	return createTestTVSeasonWithCustomURLPath("/tv/:slug/season/:seasonNumber")
}

func createTestTVSeasonWithCustomURLPath(urlPath string) (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      urlPath,
		TemplatePath: "season/item.jet",
		DataSource:   "TVSeason",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			TVSeasons: models.TVSeasonCollection{
				models.TVSeason{
					SeasonNumber: 2,
					Overview:     "# Season Title",
					ShowInfo: &models.TVShow{
						ID:        123,
						Title:     "Breaking Bad",
						TitleSlug: "breaking-bad",
						Overview:  "# Show Title",
					},
					Slug: "/tv/123/season/2",
					Episodes: []models.TVEpisode{
						models.TVEpisode{
							Overview: "# Episode Title",
						},
					},
				},
			},
		},
	}

	return ctx, r
}

func TestTVSeasonGetRouteForSlug(t *testing.T) {
	var tvSeasonDS TVSeasonDataSource

	ctx, _ := createTestTVSeason()

	route := tvSeasonDS.GetRouteForSlug(ctx, "/tv/123/season/2")

	assert.Equal(t, "/fr/tv/breaking-bad/season/2", route)
}

func TestTVSeasonIsSlugMatch(t *testing.T) {
	var tvSeasonDS TVSeasonDataSource
	assert.True(t, tvSeasonDS.IsSlugMatch("/tv/123/season/3"))
	assert.False(t, tvSeasonDS.IsSlugMatch("/tv/123"))
}

func TestTVSeasonGetRouteForMissingSlug(t *testing.T) {
	var tvSeasonDS TVSeasonDataSource

	ctx, _ := createTestTVSeason()

	route := tvSeasonDS.GetRouteForSlug(ctx, "/tv/999/season/1")

	assert.Equal(t, "ERR(/tv/999/season/1)", route)
}

func TestTVSeasonGetRouteForInvalidSlug(t *testing.T) {
	var tvSeasonDS TVSeasonDataSource

	ctx, _ := createTestTVSeason()

	route := tvSeasonDS.GetRouteForSlug(ctx, "/tv/a")

	assert.Equal(t, "ERR(/tv/a)", route)
}

func TestTVSeasonGetRouteWithShowIDForSlug(t *testing.T) {
	var tvSeasonDS TVSeasonDataSource

	ctx, _ := createTestTVSeasonWithCustomURLPath("/tv/:showID/season/:seasonNumber.html")

	route := tvSeasonDS.GetRouteForSlug(ctx, "/tv/123/season/2")

	assert.Equal(t, "/fr/tv/123/season/2.html", route)
}

func TestRenderTVSeason(t *testing.T) {
	var ds TVSeasonDataSource

	renderer := &test.MockRenderer{}

	ctx, _ := createTestTVSeason()

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/fr/tv/breaking-bad/season/2", renderer.FilePath)
	assert.Equal(t, "season/item.jet", renderer.TemplatePath)
}

func TestPartialRenderTVSeason(t *testing.T) {
	var ds TVSeasonDataSource

	renderer := &test.MockRenderer{}

	ctx, r := createTestTVSeason()
	r.PartialTemplatePath = "/season/partial.jet"
	r.PartialURLPath = "/partials/tv/:showID/season/:seasonNumber.html"

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/fr/partials/tv/123/season/2.html", renderer.FilePath)
	assert.Equal(t, "/season/partial.jet", renderer.TemplatePath)
}

func TestContentTransforms(t *testing.T) {
	var ds TVSeasonDataSource

	renderer := &test.MockRenderer{}

	ctx, _ := createTestTVSeason()

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")

	season, _ := renderer.Data["tvseason"].Elem().Interface().(models.TVSeason)

	assert.Equal(t, season.Overview, "<h1>Season Title</h1>\n")
	assert.Equal(t, season.ShowInfo.Overview, "<h1>Show Title</h1>\n")
	assert.Equal(t, season.Episodes[0].Overview, "<h1>Episode Title</h1>\n")
}
