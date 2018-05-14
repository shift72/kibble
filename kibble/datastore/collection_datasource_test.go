package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indiereign/shift72-kibble/kibble/models"
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
