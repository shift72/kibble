package datastore

import (
	"strings"
	"testing"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/test"
	"github.com/nicksnyder/go-i18n/i18n"
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

	ctx := models.RenderContext{
		Route: r,
	}

	fds := &FilmDataSource{
	//TODO: add source data
	}
	fds.Iterator(ctx, renderer1)

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

	ctx := models.RenderContext{
		Route: r,
	}

	fds := &FilmDataSource{}
	fds.Iterator(ctx, renderer1)

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

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
	}

	routeRegistry := models.NewRouteRegistry()
	routeRegistry.Add(r)

	view := models.CreateTemplateView(&routeRegistry, i18n.IdentityTfunc(), ctx)

	tem, _ := view.LoadTemplate("", "{{ routeToSlug(film.Slug, \"filmItem\") }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	fds := &FilmDataSource{}
	fds.Iterator(ctx, renderer)

	if renderer.Result.Output() != "/fr/film-special/2" {
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
	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
	}

	routeRegistry := models.NewRouteRegistry()
	routeRegistry.Add(r)

	view := models.CreateTemplateView(&routeRegistry, i18n.IdentityTfunc(), ctx)

	tem, _ := view.LoadTemplate("", "{{ routeTo(film, \"filmItem\") }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	fds := &FilmDataSource{}
	fds.Iterator(ctx, renderer)

	if renderer.Result.Output() != "/fr/film-special/2" {
		t.Errorf("Unexpected output. `%s`", renderer.Result.Output())
	}
}

func TestTransLanguage(t *testing.T) {

	Init()

	r := &models.Route{
		Name:         "filmItem",
		URLPath:      "/film-special/:filmID",
		TemplatePath: "film/item.jet",
		DataSource:   "Film",
	}

	ctx := models.RenderContext{
		Route: r,
	}

	routeRegistry := models.NewRouteRegistry()
	routeRegistry.Add(r)

	view := models.CreateTemplateView(&routeRegistry, i18n.IdentityTfunc(), ctx)

	tem, _ := view.LoadTemplate("", "MSG {{ i18n(\"settings_title\") }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	fds := &FilmDataSource{}
	fds.Iterator(ctx, renderer)

	if renderer.Result.Output() != "MSG settings_title" {
		t.Errorf("Unexpected output. `%s`", renderer.Result.Output())
	}
}
