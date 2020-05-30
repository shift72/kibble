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

	"github.com/stretchr/testify/assert"

	"kibble/models"
	"kibble/test"
)

func createTestCollection() (models.RenderContext, *models.Route) {
	return createTestCollectionWithCustomURLPath("/collection/:slug")
}

func createTestCollectionWithCustomURLPath(urlPath string) (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      urlPath,
		TemplatePath: "collection/:type.jet",
		DataSource:   "Collection",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			Collections: models.CollectionCollection{
				models.Collection{
					ID:        123,
					Slug:      "/collection/123",
					TitleSlug: "all-the-best-films",
				},
				models.Collection{
					ID:        111,
					Slug:      "/collection/111",
					TitleSlug: "movies-to-help-with-constipation",
				},
			},
		},
	}

	return ctx, r
}

func TestCollectionGetRouteForSlug(t *testing.T) {
	var collectionDS CollectionDataSource

	ctx, _ := createTestCollection()

	route := collectionDS.GetRouteForSlug(ctx, "/collection/123")

	if route != "/fr/collection/all-the-best-films" {
		t.Errorf("expected /fr/collection/all-the-best-films got %s", route)
	}
}

func TestCollectionIsSlugMatch(t *testing.T) {
	var collectionDS CollectionDataSource

	if !collectionDS.IsSlugMatch("/collection/123") {
		t.Errorf("expected /collection/123 to match")
	}

	if !collectionDS.IsSlugMatch("/feature/123") {
		t.Errorf("expected /feature/123 to match")
	}
}

func TestCollectionGetRouteForMissingSlug(t *testing.T) {
	var collectionDS CollectionDataSource

	ctx, _ := createTestCollection()

	route := collectionDS.GetRouteForSlug(ctx, "/collection/999")

	if route != "ERR(/collection/999)" {
		t.Errorf("expected ERR(/collection/999) got %s", route)
	}
}

func TestCollectionGetRouteForInvalidSlug(t *testing.T) {
	var collectionDS CollectionDataSource

	ctx, _ := createTestCollection()

	route := collectionDS.GetRouteForSlug(ctx, "/collection/a")

	if route != "ERR(/collection/a)" {
		t.Errorf("expected ERR(/collection/a) got %s", route)
	}
}

func TestCollectionGetRouteWithIDForSlug(t *testing.T) {
	var collectionDS CollectionDataSource

	ctx, _ := createTestCollectionWithCustomURLPath("/collection/:collectionID.html")

	route := collectionDS.GetRouteForSlug(ctx, "/collection/111")

	assert.Equal(t, "/fr/collection/111.html", route)
}

func TestRenderCollection(t *testing.T) {
	var ds CollectionDataSource

	ctx, _ := createTestCollection()
	renderer := &test.MockRenderer{}

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, renderer.FilePath, "/fr/collection/movies-to-help-with-constipation")
	assert.Equal(t, "collection/:type.jet", renderer.TemplatePath)
}

func TestPartialRenderCollection(t *testing.T) {
	var ds CollectionDataSource

	ctx, _ := createTestCollection()
	ctx.Route.PartialTemplatePath = "/collection/partial.jet"
	ctx.Route.PartialURLPath = "/partials/collection/:collectionID.html"

	renderer := &test.MockRenderer{}

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, renderer.FilePath, "/fr/partials/collection/111.html")
	assert.Equal(t, "/collection/partial.jet", renderer.TemplatePath)
}
