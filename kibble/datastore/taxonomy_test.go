package datastore

import (
	"testing"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/test"
)

func TestTaxonomyDataStore(t *testing.T) {

	view := jet.NewHTMLSet("../templates/")
	view.AddGlobal("version", "v1.1.145")

	renderer1 := &test.InMemoryRenderer{View: view}

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

	fds := &FilmDataSource{}
	fds.Iterator(ctx, renderer1)

	if renderer1.ErrorCount() != 0 {
		t.Error("Unexpected errors")
	}
}
