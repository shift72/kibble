package datastore

import (
	"testing"

	"kibble/models"
	"kibble/test"
	"github.com/stretchr/testify/assert"
)

func setupTestTVEpisode(urlPath string) (models.RenderContext, *models.Route) {
	r := &models.Route{
		URLPath:      urlPath,
		TemplatePath: "episode/item.jet",
		DataSource:   "TVEpisode",
	}

	season := models.TVSeason{
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
				EpisodeNumber: 10,
				Slug:          "/tv/123/season/2/episode/10",
				Overview:      "# Episode Title",
			},
		},
	}

	season.Episodes[0].Season = &season

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			TVSeasons:  models.TVSeasonCollection{season},
			TVEpisodes: models.TVEpisodeCollection{season.Episodes[0]},
		},
	}

	return ctx, r
}

func TestTVEpisodeGetRouteForSlug(t *testing.T) {
	var ds TVEpisodeDataSource
	ctx, _ := setupTestTVEpisode("/tv/:slug/season/:seasonNumber/e/:episodeNumber")
	route := ds.GetRouteForSlug(ctx, "/tv/123/season/2/episode/10")

	assert.Equal(t, "/fr/tv/breaking-bad/season/2/e/10", route)
}

func TestTVEpisodeIsSlugMatch(t *testing.T) {
	var ds TVEpisodeDataSource
	assert.False(t, ds.IsSlugMatch("/tv/123/season/3"), "season slug")
	assert.False(t, ds.IsSlugMatch("/tv/123"), "show slug")
	assert.True(t, ds.IsSlugMatch("/tv/1/season/1/episode/13"), "episode slug")
}

func TestTVEpisodeGetRouteForMissingSlug(t *testing.T) {
	var ds TVEpisodeDataSource
	ctx, _ := setupTestTVEpisode("")
	route := ds.GetRouteForSlug(ctx, "/tv/999/season/1/episode/1")

	assert.Equal(t, "ERR(/tv/999/season/1/episode/1)", route)
}

func TestTVEpisodeGetRouteForInvalidSlug(t *testing.T) {
	var ds TVEpisodeDataSource
	ctx, _ := setupTestTVEpisode("")
	route := ds.GetRouteForSlug(ctx, "/tv/a")

	assert.Equal(t, "ERR(/tv/a)", route)
}

func TestTVEpisodeGetRouteWithShowIDForSlug(t *testing.T) {
	var ds TVEpisodeDataSource
	ctx, _ := setupTestTVEpisode("/tv/:showID/season/:seasonNumber/episode/:episodeNumber.html")
	route := ds.GetRouteForSlug(ctx, "/tv/123/season/2/episode/10")

	assert.Equal(t, "/fr/tv/123/season/2/episode/10.html", route)
}

func TestRenderTVEpisode(t *testing.T) {
	var ds TVEpisodeDataSource
	renderer := &test.MockRenderer{}
	ctx, _ := setupTestTVEpisode("/tv/:slug/s:seasonNumber_e:episodeNumber.html")

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/fr/tv/breaking-bad/s2_e10.html", renderer.FilePath)
	assert.Equal(t, "episode/item.jet", renderer.TemplatePath)
}

func TestPartialRenderTVEpisode(t *testing.T) {
	var ds TVEpisodeDataSource
	renderer := &test.MockRenderer{}
	ctx, r := setupTestTVEpisode("")
	r.PartialTemplatePath = "/episode/partial.jet"
	r.PartialURLPath = "/partials/tv/:showID/season/:seasonNumber/e/:episodeNumber.html"

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/fr/partials/tv/123/season/2/e/10.html", renderer.FilePath)
	assert.Equal(t, "/episode/partial.jet", renderer.TemplatePath)
}
