package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

func createTestTVSeason() (models.RenderContext, *models.Route) {
	return createTestTVSeasonWithCustomURLPath("/tv/:slug/season/:seasonNumber")
}

func createTestTVSeasonWithCustomURLPath(urlPath string) (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      urlPath,
		TemplatePath: "tv/:type.jet",
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
