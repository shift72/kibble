package datastore

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/utils"
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

	for _, f := range ctx.Site.Films {
		filePath := ds.GetRouteForEntity(ctx, &f)
		data.Set("film", transformFilm(f))
		data.Set("site", ctx.Site)
		errCount += renderer.Render(ctx.Route, filePath, data)
	}

	return
}

// GetRouteForEntity - get the route
func (ds *FilmDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Film)
	if ok {
		return ds.GetRouteForSlug(ctx, o.Slug)
	}
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *FilmDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	filmID, err := utils.ParseIntFromSlug(slug, 2)
	if err != nil {
		return fmt.Sprintf("ERR(%s)", slug)
	}

	// supports having /:filmID in a path and /:slug
	if strings.Contains(ctx.Route.URLPath, ":filmID") {
		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":filmID", strconv.Itoa(filmID), 1)
	} else if strings.Contains(ctx.Route.URLPath, ":slug") {
		film, err := ctx.Site.Films.FindFilmByID(filmID)
		if err != nil {
			return fmt.Sprintf("ERR(%s)", slug)
		}

		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":slug", film.TitleSlug, 1)
	}
	return models.ErrDataSource
}

// IsSlugMatch - checks if the slug is a match
func (ds *FilmDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/film/")
}
