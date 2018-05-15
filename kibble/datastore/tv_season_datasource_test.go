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
					ShowInfo: &models.TVShow{
						ID:        123,
						Title:     "Breaking Bad",
						TitleSlug: "breaking-bad",
					},
					Slug: "/tv/123/season/2",
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
