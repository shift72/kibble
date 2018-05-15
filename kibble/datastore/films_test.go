package datastore

import (
	"strings"
	"testing"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/test"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/stretchr/testify/assert"
)

func createTestFilm() (models.RenderContext, *models.Route) {
	r := &models.Route{
		URLPath:      "/film/:filmID/:slug",
		TemplatePath: "film/item.jet",
		DataSource:   "Film",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "",
		Site: &models.Site{
			Films: models.FilmCollection{
				models.Film{
					ID:        123,
					Slug:      "/film/123",
					TitleSlug: "the-big-lebowski",
				},
			},
		},
	}

	return ctx, r
}
func TestApplyContentTransforms(t *testing.T) {
	var s = `# header`

	html := models.ApplyContentTransforms(s)

	if strings.TrimSpace(string(html)) != "<h1>header</h1>" {
		t.Errorf("Expectation failed. %s", html)
	}
}

func TestTransformFilm(t *testing.T) {

	f := &models.Film{
		Overview: "# header",
	}

	tf := transformFilm(*f)

	if tf.Overview == f.Overview {
		t.Error("Expect not side effects")
	}

	if tf.Overview != "<h1>header</h1>\n" {
		t.Errorf("Expect markdown to be applied. %s\n", tf.Overview)
	}
}

func TestFilmDataStore(t *testing.T) {

	view := jet.NewHTMLSet("../templates/")
	view.AddGlobal("version", "v1.1.145")

	ctx, _ := createTestFilm()
	renderer1 := &test.InMemoryRenderer{View: view}

	var fds FilmDataSource
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
		Route:       r,
		RoutePrefix: "",
		Site: &models.Site{
			Films: models.FilmCollection{
				models.Film{
					ID:        123,
					Slug:      "/film/123",
					TitleSlug: "the-big-lebowski",
				},
			},
		},
	}

	var fds FilmDataSource
	fds.Iterator(ctx, renderer1)

	if renderer1.Result.Output() != "v1.1.145" {
		t.Error("Unexpected output")
	}
}

func TestRenderingSlug(t *testing.T) {

	Init()

	r := &models.Route{
		Name:         "filmItem",
		URLPath:      "/film-special/:slug",
		TemplatePath: "film/item.jet",
		DataSource:   "Film",
	}

	cfg := models.Config{
		Routes: []models.Route{*r},
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			Films: models.FilmCollection{
				models.Film{
					ID:        123,
					Slug:      "/film/123",
					TitleSlug: "the-big-lebowski",
				},
			},
		},
	}

	routeRegistry := models.NewRouteRegistryFromConfig(&cfg)

	view := models.CreateTemplateView(routeRegistry, i18n.IdentityTfunc(), ctx, "./templates")

	tem, _ := view.LoadTemplate("", "{{ routeToSlug(film.Slug, \"filmItem\") }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	var fds FilmDataSource
	fds.Iterator(ctx, renderer)

	if renderer.Result.Output() != "/fr/film-special/the-big-lebowski" {
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
		Site: &models.Site{
			Films: models.FilmCollection{
				models.Film{
					ID:        123,
					Slug:      "/film/123",
					TitleSlug: "the-big-lebowski",
				},
			},
		},
	}

	cfg := models.Config{
		Routes: []models.Route{*r},
	}
	routeRegistry := models.NewRouteRegistryFromConfig(&cfg)

	view := models.CreateTemplateView(routeRegistry, i18n.IdentityTfunc(), ctx, "./templates")

	tem, _ := view.LoadTemplate("", "{{ routeTo(film, \"filmItem\") }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	var fds FilmDataSource
	fds.Iterator(ctx, renderer)

	if renderer.Result.Output() != "/fr/film-special/123" {
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
		Site: &models.Site{
			Films: models.FilmCollection{
				models.Film{
					ID:        123,
					Slug:      "/film/123",
					TitleSlug: "the-big-lebowski",
				},
			},
		},
	}

	cfg := models.Config{
		Routes: []models.Route{*r},
	}

	routeRegistry := models.NewRouteRegistryFromConfig(&cfg)

	view := models.CreateTemplateView(routeRegistry, i18n.IdentityTfunc(), ctx, "./templates")

	tem, _ := view.LoadTemplate("", "MSG {{ i18n(\"settings_title\") }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	var fds FilmDataSource
	fds.Iterator(ctx, renderer)

	if renderer.Result.Output() != "MSG settings_title" {
		t.Errorf("Unexpected output. `%s`", renderer.Result.Output())
	}
}

func TestRenderFilm(t *testing.T) {
	ctx, _ := createTestFilm()
	renderer := &test.MockRenderer{}

	var ds FilmDataSource
	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/film/123/the-big-lebowski", renderer.FilePath)
	assert.Equal(t, "film/item.jet", renderer.TemplatePath)
}

func TestRenderPartialFilm(t *testing.T) {
	ctx, r := createTestFilm()
	r.PartialTemplatePath = "/film/partial.jet"
	r.PartialURLPath = "/partials/film/:filmID.html"

	renderer := &test.MockRenderer{}

	var ds FilmDataSource
	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/partials/film/123.html", renderer.FilePath)
	assert.Equal(t, "/film/partial.jet", renderer.TemplatePath)
}
