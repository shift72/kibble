package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

func createTestTVShow() (models.RenderContext, *models.Route) {
	return createTestTVShowWithCustomURLPath("/tv/:slug")
}

func createTestTVShowWithCustomURLPath(urlPath string) (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      urlPath,
		TemplatePath: "tv/:type.jet",
		DataSource:   "TVShow",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			TVShows: models.TVShowCollection{
				models.TVShow{
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
