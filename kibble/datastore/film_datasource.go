package datastore

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/pressly/chi"
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

// Query - return a single film
func (ds *FilmDataSource) Query(ctx models.RenderContext, req *http.Request) (jet.VarMap, error) {

	var f *models.Film

	if strings.Contains(ctx.Route.URLPath, ":filmID") {

		filmID, err := strconv.Atoi(chi.URLParam(req, "filmID"))
		if err != nil {
			return nil, err
		}

		f, err = ctx.Site.Films.FindFilmByID(filmID)
		if err != nil || f == nil {
			return nil, err
		}

	} else if strings.Contains(ctx.Route.URLPath, ":slug") {
		slug := chi.URLParam(req, "slug")
		f, _ = ctx.Site.Films.FindFilmBySlug(slug)
	}

	if f == nil {
		//TODO: indicate a 404 error
		return nil, errors.New("Not found")
	}

	vars := make(jet.VarMap)
	vars.Set("film", transformFilm(*f))
	vars.Set("site", ctx.Site)
	return vars, nil
}

// Iterator - loop over each film
func (ds *FilmDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	data := make(jet.VarMap)

	for _, f := range ctx.Site.Films {

		filePath := ds.GetRouteForEntity(ctx, &f)

		c := transformFilm(f)

		data.Set("film", c)
		data.Set("site", ctx.Site)
		renderer.Render(ctx.Route, filePath, data)
	}
}

// GetRouteForEntity - get the route
func (ds *FilmDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	o, ok := entity.(*models.Film)
	if ok {
		return ds.GetRouteForSlug(ctx, o.Slug)
	}
	return models.DataSourceError
}

// GetRouteForSlug - get the route
func (ds *FilmDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	//TODO: parse slug
	p := strings.Split(slug, "/")
	if strings.Contains(ctx.Route.URLPath, ":filmID") {
		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":filmID", p[2], 1)
	} else if strings.Contains(ctx.Route.URLPath, ":slug") {
		filmID, _ := strconv.Atoi(p[2])
		film, _ := ctx.Site.Films.FindFilmByID(filmID)
		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":slug", film.TitleSlug, 1)
	}
	return models.DataSourceError
}

// IsSlugMatch - checks if the slug is a match
func (ds *FilmDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/film/")
}
