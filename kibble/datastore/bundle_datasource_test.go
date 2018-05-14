package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/test"
)

func createTestContext() (models.RenderContext, *models.Route) {

	r := &models.Route{
		URLPath:      "/bundle/:bundleID/:slug",
		TemplatePath: "bundle/item.jet",
		DataSource:   "Bundle",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "",
		Site: &models.Site{
			Bundles: models.BundleCollection{
				models.Bundle{
					ID:        111,
					TitleSlug: "marks-big-saggy-bundle",
					Title:     "Marks Big Saggy Bundle",
				},
			},
		},
	}

	return ctx, r
}

func TestRenderForBundle(t *testing.T) {
	var ds BundleDataSource

	renderer := &test.MockRenderer{}

	ctx, _ := createTestContext()

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/bundle/111/marks-big-saggy-bundle", renderer.FilePath)
	assert.Equal(t, "bundle/item.jet", renderer.Route.TemplatePath)
}

func TestPartialRenderForBundle(t *testing.T) {
	var ds BundleDataSource

	renderer := &test.MockRenderer{}

	ctx, r := createTestContext()
	r.PartialTemplatePath = "/bundle/partial.jet"
	r.PartialURLPath = "/partials/bundle/:bundleID.html"

	ds.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled, "renderer.RenderCalled")
	assert.Equal(t, "/partials/bundle/111.html", renderer.FilePath)
	assert.Equal(t, "/bundle/partial.jet", renderer.Route.TemplatePath)
}
