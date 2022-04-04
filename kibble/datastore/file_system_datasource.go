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
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/CloudyKit/jet/v6"
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

	// check if we should render these static files
	if ctx.Route.DefaultLanguageOnly && ctx.RoutePrefix != "" {
		return
	}

	data := make(jet.VarMap)
	data.Set("site", ctx.Site)

	dirPath := filepath.Join(ctx.Site.SiteConfig.SiteRootPath, ctx.Route.TemplatePath)
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		log.Warningf("FileSystem data source could not find path: %s ", dirPath)
		return
	}

	searchPath := filepath.Join(ctx.Site.SiteConfig.SiteRootPath, ctx.Route.TemplatePath, "*.jet")

	files, _ := filepath.Glob(searchPath)
	for _, f := range files {
		relativeFilePath := strings.Replace(f, ctx.Site.SiteConfig.SiteRootPath, "", 1)
		urlPath := path.Join(ctx.RoutePrefix, strings.Replace(relativeFilePath, ".jet", "", 1))
		errCount += renderer.Render(relativeFilePath, urlPath, data)
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
