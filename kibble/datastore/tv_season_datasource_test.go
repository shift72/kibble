package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

func createTestTVSeason() (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      "/tv/:slug/season/:seasonNumber",
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
					ShowInfo: models.TVShow{
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

	if !tvSeasonDS.IsSlugMatch("/tv/123/season/2") {
		t.Errorf("expected /tv/123/season/2 to match")
	}

	if !tvSeasonDS.IsSlugMatch("/tv/123") {
		t.Errorf("expected /tv/123 to match")
	}
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
