package datastore

import (
	"testing"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/test"
)

func createTestContextHomepage() (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      "/page/:slug",
		TemplatePath: "page/:type.jet",
		DataSource:   "Page",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site: &models.Site{
			Pages: models.PageCollection{
				models.Page{
					ID:       123,
					Slug:     "homepage-slug",
					PageType: "homepage",
				},
			},
		},
	}

	return ctx, r
}

func createTestContextCurated() (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      "/page/:slug",
		TemplatePath: "page/:type.jet",
		DataSource:   "Page",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "",
		Site: &models.Site{
			Pages: models.PageCollection{
				models.Page{
					ID:       123,
					Slug:     "disney",
					PageType: "curated",
				},
			},
		},
	}

	return ctx, r
}

func createTestContextExternal() models.RenderContext {

	r := &models.Route{
		URLPath:      "/page/:slug",
		TemplatePath: "page/:type.jet",
		DataSource:   "Page",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "",
		Site: &models.Site{
			Pages: models.PageCollection{
				models.Page{
					ID:       1,
					Slug:     "about-us",
					PageType: "content",
				},
				models.Page{
					ID:       123,
					PageType: "external",
					URL:      "https://www.shift72.com",
				},
				models.Page{
					ID:       124,
					PageType: "external",
					URL:      "https://www.shift72.com/#!/page/about-us",
				},
			},
			SiteConfig: &models.Config{
				SiteURL: "https://www.shift72.com",
			},
		},
	}

	return ctx
}

func createTestContextExternalOnly() models.RenderContext {

	r := &models.Route{
		URLPath:      "/page/:slug",
		TemplatePath: "page/:type.jet",
		DataSource:   "Page",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "",
		Site: &models.Site{
			Pages: models.PageCollection{
				models.Page{
					ID:       123,
					PageType: "external",
					URL:      "https://www.shift72.com",
				},
			},
		},
	}

	return ctx
}

func TestHomepageTemplateType(t *testing.T) {

	pageDS := PageDataSource{}

	renderer := &test.MockRenderer{}

	ctx, _ := createTestContextHomepage()

	pageDS.Iterator(ctx, renderer)

	if !renderer.RenderCalled {
		t.Error("Expected render to be called")
	}

	if renderer.Route.TemplatePath != "page/homepage.jet" {
		t.Errorf("Expected render template to be '/page/item.jet' got %s\n", renderer.Route.TemplatePath)
	}

	if renderer.FilePath != "/fr/" {
		t.Errorf("Expected file path to be '/fr/' got %s\n", renderer.FilePath)
	}
}

func TestCuratedTemplateType(t *testing.T) {

	pageDS := PageDataSource{}

	renderer := &test.MockRenderer{}

	ctx, r := createTestContextCurated()

	pageDS.Iterator(ctx, renderer)

	if !renderer.RenderCalled {
		t.Error("Expected render to be called")
	}

	if renderer.Route.TemplatePath != "page/curated.jet" {
		t.Errorf("Expected render template to be '/page/curated.jet' got %s\n", renderer.Route.TemplatePath)
	}

	if renderer.FilePath != "/page/disney" {
		t.Errorf("Expected file path to be '/page/disney/' got %s\n", renderer.FilePath)
	}

	if r.TemplatePath != "page/:type.jet" {
		t.Errorf("Expected render template to be '/page/:type.jet' got %s\n", r.TemplatePath)
	}
}

func TestExternalTemplateType(t *testing.T) {

	pageDS := PageDataSource{}

	renderer := &test.MockRenderer{}

	ctx := createTestContextExternalOnly()

	pageDS.Iterator(ctx, renderer)

	if renderer.RenderCalled {
		t.Error("Expected render to be not be called for external pages")
	}
}

func TestGetRouteForExternalPage(t *testing.T) {
	pageDS := PageDataSource{}

	ctx := createTestContextExternal()

	route := pageDS.GetRouteForSlug(ctx, "/page/123")

	if route != "https://www.shift72.com" {
		t.Error("expected https://www.shift72.com got ", route)
	}

	route = pageDS.GetRouteForSlug(ctx, "/page/124")

	if route != "/page/about-us" {
		t.Error("expected /page/about-us ", route)
	}
}

func TestParse(t *testing.T) {
	pageDS := PageDataSource{}

	if !pageDS.IsSlugMatch("/page/2") {
		t.Error("expected /page/2")
	}

	if pageDS.IsSlugMatch("/film/2") {
		t.Error("expected /film/2 should fail")
	}
}

func TestGetEntityType(t *testing.T) {
	pageDS := PageDataSource{}

	if pageDS.GetEntityType().String() != "*models.Page" {
		t.Error("expected *models.Page")
	}
}

func TestGetRouteForSlug(t *testing.T) {
	pageDS := PageDataSource{}

	ctx, _ := createTestContextCurated()

	route := pageDS.GetRouteForSlug(ctx, "/page/123")

	if route != "/page/disney" {
		t.Error("expected /page/disney")
	}
}

func TestGetRouteForInvalidSlug(t *testing.T) {
	pageDS := PageDataSource{}

	ctx, _ := createTestContextCurated()

	route := pageDS.GetRouteForSlug(ctx, "/page/a")

	if route != "ERR(/page/a)" {
		t.Errorf("expected ERR(/page/a) got %s", route)
	}
}

func TestGetRouteForMissingSlug(t *testing.T) {
	pageDS := PageDataSource{}

	ctx, _ := createTestContextCurated()

	route := pageDS.GetRouteForSlug(ctx, "/page/999")

	if route != "ERR(/page/999)" {
		t.Errorf("expected ERR(/page/999) got %s", route)
	}
}
