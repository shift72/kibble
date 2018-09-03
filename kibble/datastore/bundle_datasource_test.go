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
	assert.Equal(t, "bundle/item.jet", renderer.TemplatePath)
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
	assert.Equal(t, "/bundle/partial.jet", renderer.TemplatePath)
}

func TestIsValid(t *testing.T) {
	var ds BundleDataSource
	assert.NoError(t, models.ValidateRouteWithDatasource("/bundle/:bundleID/:slug", &ds))
}
func TestIsNotValid(t *testing.T) {
	var ds BundleDataSource
	assert.EqualError(t,
		models.ValidateRouteWithDatasource("/bundle/:bunID/:slug", &ds),
		"Path (/bundle/:bunID/:slug) contains invalid replacement arguments. Valid arguments are (:bundleID,:slug)")
}
