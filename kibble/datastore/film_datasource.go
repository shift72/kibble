package datastore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// FilmDataSource - single film datasource
// Supports slugs in the /film/:filmID and /film/:title_slug
type FilmDataSource struct{}

// GetName - name of the datasource
func (ds *FilmDataSource) GetName() string {
	return "Film"
}

// GetEntityType - Get the entity type
func (ds *FilmDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf(&models.Film{})
}

// Iterator - loop over each film
func (ds *FilmDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	for _, f := range ctx.Site.Films {
		data.Set("film", transformFilm(f))

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
func (ds *FilmDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Film)
	if ok {
		s := strings.Replace(ctx.Route.URLPath, ":filmID", strconv.Itoa(o.ID), 1)
		s = strings.Replace(s, ":slug", o.TitleSlug, 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
}

// GetPartialRouteForEntity - get the partial route
func (ds *FilmDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Film)
	if ok {
		s := strings.Replace(ctx.Route.PartialURLPath, ":filmID", strconv.Itoa(o.ID), 1)
		s = strings.Replace(s, ":slug", o.TitleSlug, 1)
		return ctx.RoutePrefix + s
	}
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *FilmDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	film, ok := ctx.Site.Films.FindFilmBySlug(slug)
	if !ok {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	return ds.GetRouteForEntity(ctx, film)
}

// IsSlugMatch - checks if the slug is a match
func (ds *FilmDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/film/")
}
