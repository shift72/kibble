package datastore

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/pressly/chi"
)

// FilmDataSource - single film datasource
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
func (ds *FilmDataSource) Query(req *http.Request) (jet.VarMap, error) {

	filmID, err := strconv.Atoi(chi.URLParam(req, "filmID"))
	if err != nil {
		return nil, err
	}

	f, err := FindByID(filmID)
	if err != nil || f == nil {
		return nil, err
	}
	c := transformFilm(*f)

	vars := make(jet.VarMap)
	vars.Set("film", c)
	return vars, nil
}

// Iterator - loop over each film
func (ds *FilmDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) {

	films, _ := GetAllFilms()
	data := make(jet.VarMap)

	for _, f := range *films {

		filePath := ds.GetRouteForEntity(ctx, &f)

		c := transformFilm(f)

		data.Set("film", c)
		renderer.Render(ctx.Route, filePath, data)
	}
}

// GetRouteForEntity - get the route
func (ds *FilmDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {

	o, ok := entity.(*models.Film)
	if ok {
		return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":filmID", strconv.Itoa(o.ID), 1)
	}
	return "!Error"
}

// GetRouteForSlug - get the route
func (ds *FilmDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	//TODO: parse slug
	p := strings.Split(slug, "/")
	return ctx.RoutePrefix + strings.Replace(ctx.Route.URLPath, ":filmID", p[2], 1)
}

// IsSlugMatch - checks if the slug is a match
func (ds *FilmDataSource) IsSlugMatch(slug string) bool {
	return strings.HasPrefix(slug, "/film/")
}
