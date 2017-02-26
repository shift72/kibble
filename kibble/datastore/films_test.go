package datastore

import (
	"strings"
	"testing"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/test"
)

func TestApplyContentTransforms(t *testing.T) {
	var s = `# header`

	html := applyContentTransforms(s)

	if strings.TrimSpace(string(html)) != "<h1>header</h1>" {
		t.Errorf("Expectation failed. %s", html)
	}
}

func TestTransformFilm(t *testing.T) {

	f := &models.Film{
		Synopsis: "# header",
	}

	tf := transformFilm(*f)

	if tf.Synopsis == f.Synopsis {
		t.Error("Expect not side effects")
	}

	if tf.Synopsis != "<h1>header</h1>\n" {
		t.Errorf("Expect markdown to be applied. %s\n", tf.Synopsis)
	}
}

func TestFilmDataStore(t *testing.T) {

	view := jet.NewHTMLSet("../templates/")
	view.AddGlobal("version", "v1.1.145")

	renderer1 := &test.InMemoryRenderer{View: view}

	r := &models.Route{
		URLPath:      "/film/:filmID",
		TemplatePath: "film/item.jet",
		DataSource:   "Film",
	}
	fds := &FilmDataSource{
	//TODO: add source data
	}
	fds.Iterator(r, renderer1)

	if renderer1.ErrorCount() != 0 {
		t.Error("Unexpected errors")
	}
}

func TestRenderingGlobal(t *testing.T) {

	view := jet.NewHTMLSet("../templates/")
	view.AddGlobal("version", "v1.1.145")

	tem, _ := view.LoadTemplate("", "{{ version }}")

	renderer1 := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	r := &models.Route{
		URLPath:      "/film/:filmID",
		TemplatePath: "film/item.jet",
		DataSource:   "Film",
	}
	fds := &FilmDataSource{}
	fds.Iterator(r, renderer1)

	if renderer1.Result.Output() != "v1.1.145" {
		t.Error("Unexpected output")
	}
}

func TestRenderingSlug(t *testing.T) {

	Init()

	r := &models.Route{
		Name:         "filmItem",
		URLPath:      "/film-special/:filmID",
		TemplatePath: "film/item.jet",
		DataSource:   "Film",
	}

	routeRegistry := models.NewRouteRegistry()
	routeRegistry.Add(r)

	view := jet.NewHTMLSet("../templates/")
	view.AddGlobal("version", "v1.1.145")
	view.AddGlobal("routeTo", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(entity, routeName)
	})
	view.AddGlobal("routeToSlug", func(slug string, routeName string) string {
		return routeRegistry.GetRouteForSlug(slug, routeName)
	})

	tem, _ := view.LoadTemplate("", "{{ routeToSlug(film.Slug, \"filmItem\") }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	fds := &FilmDataSource{}
	fds.Iterator(r, renderer)

	if renderer.Result.Output() != "/film-special/2" {
		t.Errorf("Unexpected output. `%s`", renderer.Result.Output())
	}
}

func TestRouteToFilm(t *testing.T) {

	Init()

	r := &models.Route{
		Name:         "filmItem",
		URLPath:      "/film-special/:filmID",
		TemplatePath: "film/item.jet",
		DataSource:   "Film",
	}

	routeRegistry := models.NewRouteRegistry()
	routeRegistry.Add(r)

	view := jet.NewHTMLSet("../templates/")
	view.AddGlobal("version", "v1.1.145")
	view.AddGlobal("routeTo", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(entity, routeName)
	})
	view.AddGlobal("routeToSlug", func(slug string, routeName string) string {
		return routeRegistry.GetRouteForSlug(slug, routeName)
	})

	tem, _ := view.LoadTemplate("", "{{ routeTo(film, \"filmItem\") }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	fds := &FilmDataSource{}
	fds.Iterator(r, renderer)

	if renderer.Result.Output() != "/film-special/2" {
		t.Errorf("Unexpected output. `%s`", renderer.Result.Output())
	}
}
