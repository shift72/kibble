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
	data.Set("site", ctx.Site)

	for _, f := range ctx.Site.TVShows {
		data.Set("tvshow", transformTVShow(f))

		filePath := ds.GetRouteForEntity(ctx, &f)
		errCount += renderer.Render(ctx.Route.TemplatePath, filePath, data)

		if ctx.Route.HasPartial() {
			partialFilePath := ds.GetPartialRouteForEntity(ctx, &f)
			errCount += renderer.Render(ctx.Route.PartialTemplatePath, partialFilePath, data)
		}
	}

	return
}

// GetRouteForEntity - get the route
func (ds *TVShowDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.TVShow)
	log.Critical("%s", ok)
	if ok {
		s := strings.Replace(ctx.Route.URLPath, ":slug", o.TitleSlug, 1)
		s = strings.Replace(s, ":showID", strconv.Itoa(o.ID), 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
}

// GetPartialRouteForEntity - get the partial route
func (ds *TVShowDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.TVShow)
	if ok {
		s := strings.Replace(ctx.Route.PartialURLPath, ":slug", o.TitleSlug, 1)
		s = strings.Replace(s, ":showID", strconv.Itoa(o.ID), 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *TVShowDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {

	// supports any params: :slug, :showID
	tvShow, found := ctx.Site.TVShows.FindTVShowBySlug(slug)
	if found {
		return ds.GetRouteForEntity(ctx, tvShow)
	}
	return fmt.Sprintf("ERR(%s)", slug)
}

// IsSlugMatch - checks if the slug is a match
func (ds *TVShowDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/tv/") && !strings.Contains(slug, "/season/")
}
