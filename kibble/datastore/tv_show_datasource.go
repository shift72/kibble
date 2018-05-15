package datastore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// TVShowDataSource - single tv season datasource
// Supports slugs in the /tv/:tvID/season/:seasonID and /tv/:title_slug
type TVShowDataSource struct{}

// GetName - name of the datasource
func (ds *TVShowDataSource) GetName() string {
	return "TVShow"
}

// GetEntityType - Get the entity type
func (ds *TVShowDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.TVShow{})
}

// Iterator - loop over each film
func (ds *TVShowDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)

	for _, f := range ctx.Site.TVShows {
		filePath := ds.GetRouteForEntity(ctx, &f)

		data.Set("tvshow", transformTVShow(f))
		data.Set("site", ctx.Site)
		errCount += renderer.Render(ctx.Route.TemplatePath, filePath, data)
	}

	return
}

// GetRouteForEntity - get the route
func (ds *TVShowDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.TVShow)
	if ok {
		return ds.GetRouteForSlug(ctx, o.Slug)
	}
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *TVShowDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {

	// supports any params: :slug, :showID
	tvShow, found := ctx.Site.TVShows.FindTVShowBySlug(slug)
	if !found {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	s := strings.Replace(ctx.Route.URLPath, ":slug", tvShow.TitleSlug, 1)
	s = strings.Replace(s, ":showID", strconv.Itoa(tvShow.ID), 1)

	return ctx.RoutePrefix + s
}

// IsSlugMatch - checks if the slug is a match
func (ds *TVShowDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/tv/") && !strings.Contains(slug, "/season/")
}
