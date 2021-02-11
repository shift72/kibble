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

	"kibble/models"
	"kibble/test"

	"github.com/stretchr/testify/assert"
)

func TestRootTemplateType(t *testing.T) {

	var fsDS FileSystemDataSource

	renderer := &test.MockRenderer{}

	r := &models.Route{
		Name:                "root",
		URLPath:             "",
		TemplatePath:        "../sample_site/",
		DataSource:          "FileSystem",
		DefaultLanguageOnly: true,
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "",
		Site:        &models.Site{},
	}
	fsDS.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled)
	assert.Equal(t, "../sample_site/test.html.jet", renderer.TemplatePath)
	assert.Equal(t, "../sample_site/test.html", renderer.FilePath)
}

func TestRootTemplateWithRoutePrefixType(t *testing.T) {

	var fsDS FileSystemDataSource

	renderer := &test.MockRenderer{}

	r := &models.Route{
		Name:                "root",
		URLPath:             "",
		TemplatePath:        ".",
		DataSource:          "FileSystem",
		DefaultLanguageOnly: true,
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr",
		Site:        &models.Site{},
	}

	fsDS.Iterator(ctx, renderer)

	assert.False(t, renderer.RenderCalled)
}

func TestTemplateWithRoutePrefixType(t *testing.T) {

	var fsDS FileSystemDataSource

	renderer := &test.MockRenderer{}

	r := &models.Route{
		Name:         "static",
		URLPath:      "",
		TemplatePath: "../sample_site/help/",
		DataSource:   "FileSystem",
	}

	ctx := models.RenderContext{
		Route:       r,
		RoutePrefix: "/fr/blah", // blah added because we have used the relative .. pathing above
		Site:        &models.Site{},
	}

	fsDS.Iterator(ctx, renderer)

	assert.True(t, renderer.RenderCalled)
	assert.Equal(t, "../sample_site/help/example.jet", renderer.TemplatePath)
	assert.Equal(t, "/fr/sample_site/help/example", renderer.FilePath)
}
