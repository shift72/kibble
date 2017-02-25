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
