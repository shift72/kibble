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
	"fmt"
	"kibble/models"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/CloudyKit/jet"
)

var fileSystemArgs = []models.RouteArgument{}

// FileSystemDataSource - represents a directory of .jet templates
type FileSystemDataSource struct {
}

// GetName - name of the datasource
func (ds *FileSystemDataSource) GetName() string {
	return "FileSystem"
}

// GetRouteArguments returns the available route arguments
func (ds *FileSystemDataSource) GetRouteArguments() []models.RouteArgument {
	return fileSystemArgs
}

// GetEntityType - Get the entity type
func (ds *FileSystemDataSource) GetEntityType() reflect.Type {
	return reflect.TypeOf("")
}

// Iterator - loop over each film
func (ds *FileSystemDataSource) Iterator(ctx models.RenderContext, renderer models.Renderer) (errCount int) {

	// special case - only render the implicit root route if this is the default language context
	// these files will not be rendered for every language
	if ctx.Route.Name == "root" && ctx.RoutePrefix != "" {
		return
	}

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	files, _ := filepath.Glob(filepath.Join(ctx.Route.TemplatePath, "*.jet"))
	for _, f := range files {
		filePath := path.Join(ctx.RoutePrefix, strings.Replace(f, ".jet", "", 1))
		errCount += renderer.Render(f, filePath, data)
	}

	return
}

// GetRouteForEntity - get the route
func (ds *FileSystemDataSource) GetRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ErrDataSource
}

// GetPartialRouteForEntity - get the partial route
func (ds *FileSystemDataSource) GetPartialRouteForEntity(ctx models.RenderContext, entity interface{}) string {
	return models.ErrDataSource
}

// GetRouteForSlug - get the route
func (ds *FileSystemDataSource) GetRouteForSlug(ctx models.RenderContext, slug string) string {
	return fmt.Sprintf("ERR(%s)", slug)
}

// IsSlugMatch - checks if the slug is a match
func (ds *FileSystemDataSource) IsSlugMatch(slug string) bool {
	return false
}

// IsValid checks for any validation errors
func (ds *FileSystemDataSource) IsValid(route *models.Route) error {
	return nil
}
