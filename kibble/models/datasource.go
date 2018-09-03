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

package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// ErrDataSource is returned when no datasource could be found
	ErrDataSource = "!Error"
	// ErrDataSourceMissing raises when the datasource could not be found
	ErrDataSourceMissing = errors.New("Missing")
)

var store map[string]DataSource

// DataSource provides a set of data for querying and iterating over
type DataSource interface {
	GetName() string
	GetEntityType() reflect.Type
	Iterator(ctx RenderContext, renderer Renderer) (errorCount int)
	IsSlugMatch(slug string) bool
	GetRouteForEntity(ctx RenderContext, entity interface{}) string
	GetRouteForSlug(ctx RenderContext, slug string) string
	GetRouteArguments() []RouteArgument
}

// AddDataSource - register a datasource
func AddDataSource(ds DataSource) {

	if store == nil {
		store = make(map[string]DataSource)
	}

	store[ds.GetName()] = ds
}

// FindDataSource - find the data source by name
func FindDataSource(name string) DataSource {
	return store[name]
}

// GetDataSources returns all registered data sources
func GetDataSources() map[string]DataSource {
	return store
}

// RouteArgument represents an argument that a route can have
type RouteArgument struct {
	Name        string
	Description string
	GetValue    func(obj interface{}) string
}

// ReplaceURLArgumentsWithEntityValues replaces expected arguments with
func ReplaceURLArgumentsWithEntityValues(routePrefix string, urlPath string, args []RouteArgument, entity interface{}) string {

	for i := 0; i < len(args); i++ {
		value := args[i].GetValue(entity)
		if value == ErrDataSource {
			return value
		}

		urlPath = strings.Replace(urlPath, args[i].Name, args[i].GetValue(entity), 1)
	}

	return routePrefix + urlPath
}

// ValidateRouteWithDatasource will check the urlPath is valid for a data source
func ValidateRouteWithDatasource(urlPath string, ds DataSource) error {

	parts := strings.FieldsFunc(urlPath, urlSplit)

	for i := range parts {
		if strings.HasPrefix(parts[i], ":") {
			found := false
			for _, a := range ds.GetRouteArguments() {
				if parts[i] == a.Name {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("Path (%s) contains invalid replacement arguments. %s", urlPath, validArguments(ds))
			}
		}
	}
	return nil
}

func urlSplit(r rune) bool {
	return r == '/' || r == '.' || r == '-'
}

func validArguments(ds DataSource) string {
	validArgs := ""
	for _, a := range ds.GetRouteArguments() {
		if validArgs == "" {
			validArgs = a.Name
		} else {
			validArgs = validArgs + "," + a.Name
		}
	}
	return "Valid arguments are (" + validArgs + ")"
}
