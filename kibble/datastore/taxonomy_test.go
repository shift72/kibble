//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

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

	var fds FilmDataSource
	fds.Iterator(ctx, renderer1)

	if renderer1.ErrorCount() != 0 {
		t.Error("Unexpected errors")
	}
}
