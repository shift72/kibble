package datastore

import (
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/test"
	"github.com/nicksnyder/go-i18n/i18n"
)

func TestRenderingShowSlug(t *testing.T) {

	Init()

	r := &models.Route{
		Name:         "tvShowItem",
		URLPath:      "/tv-show/:slug",
		TemplatePath: "tvshow/item.jet",
		DataSource:   "TVShow",
	}

	cfg := models.Config{
		Routes: []models.Route{*r},
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			TVShows: models.TVShowCollection{
				models.TVShow{
					ID:        123,
					Slug:      "/tv/123",
					TitleSlug: "the-big-lebowski",
				},
			},
		},
	}

	routeRegistry := models.NewRouteRegistryFromConfig(&cfg)

	view := models.CreateTemplateView(routeRegistry, i18n.IdentityTfunc(), ctx, "./templates")

	tem, _ := view.LoadTemplate("", "{{ routeToSlug(site.TVShows[0].Slug) }}")

	renderer := &test.InMemoryTemplateRenderer{
		View:     view,
		Template: tem,
	}

	var ds TVShowDataSource
	ds.Iterator(ctx, renderer)

	if renderer.Result.Output() != "/fr/tv-show/the-big-lebowski" {
		t.Errorf("Unexpected output. `%s`", renderer.Result.Output())
	}
}
