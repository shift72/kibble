package datastore

import (
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
)

func createTestCollection() (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      "/collection/:slug",
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
			},
		},
	}

	return ctx, r
}

func TestCollectionGetRouteForSlug(t *testing.T) {
	collectionDS := CollectionDataSource{}

	ctx, _ := createTestCollection()

	route := collectionDS.GetRouteForSlug(ctx, "/collection/123")

	if route != "/fr/collection/all-the-best-films" {
		t.Errorf("expected /fr/collection/all-the-best-films got %s", route)
	}
}

func TestCollectionGetRouteForMissingSlug(t *testing.T) {
	collectionDS := CollectionDataSource{}

	ctx, _ := createTestCollection()

	route := collectionDS.GetRouteForSlug(ctx, "/collection/999")

	if route != "ERR(/collection/999)" {
		t.Errorf("expected ERR(/collection/999) got %s", route)
	}
}

func TestCollectionGetRouteForInvalidSlug(t *testing.T) {
	collectionDS := CollectionDataSource{}

	ctx, _ := createTestCollection()

	route := collectionDS.GetRouteForSlug(ctx, "/collection/a")

	if route != "ERR(/collection/a)" {
		t.Errorf("expected ERR(/collection/a) got %s", route)
	}
}
