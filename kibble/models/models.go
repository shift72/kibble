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
	"path"
	"reflect"

	"github.com/CloudyKit/jet"
)

// Pagination describes a single page of results
type Pagination struct {
	Index       int
	Size        int
	Total       int
	PreviousURL string
	NextURL     string
}

// DataSource provides a set of data for querying and iterating over
type DataSource interface {
	GetName() string
	GetEntityType() reflect.Type
	Iterator(ctx RenderContext, renderer Renderer) (errorCount int)
	IsSlugMatch(slug string) bool
	GetRouteForEntity(ctx RenderContext, entity interface{}) string
	GetRouteForSlug(ctx RenderContext, slug string) string
	//TODO: ValidateRoute - check that the route contains valid tokens
}

// RenderContext - Context passed during rendering / serving
type RenderContext struct {
	Route       *Route
	RoutePrefix string
	Site        *Site
	Language    *Language
}

// Renderer - rendering implementation
type Renderer interface {
	Initialise()
	Render(templatePath string, filePath string, data jet.VarMap) (errorCount int)
}

// Config - template configuration
// NOTE: Don't use `SiteRootPath directly`, use `Config.SourcePath()` instead.
type Config struct {
	DefaultLanguage string            `json:"defaultLanguage"`
	Languages       map[string]string `json:"languages"`
	Routes          []Route           `json:"routes"`
	SiteURL         string            `json:"siteUrl"`
	BuilderVersion  string            `json:"builderVersion"`
	Version         string            `json:"version"`
	Name            string            `json:"Name"`
	Private         PrivateConfig     `json:"-"`
	DisableCache    bool              `json:"-"`
	RunAsAdmin      bool              `json:"-"`
	SkipLogin       bool              `json:"-"`
	SiteRootPath    string            `json:"siteRootPath"`
	LiveReload      LiveReloadConfig  `json:"liveReload"`
}

// LiveReloadConfig - configuration options for the live_reloader
type LiveReloadConfig struct {
	IgnoredPaths []string `json:"ignoredPaths"`
}

// PrivateConfig - config loaded from
type PrivateConfig struct {
	APIKey string `json:"apikey"`
}

// BuildPath returns the build path for current config
func (cfg *Config) BuildPath() string {
	if cfg.RunAsAdmin {
		return path.Join(".kibble", "build-admin")
	}
	return path.Join(".kibble", "build")
}

// FileRootPath returns the path to be used for copying
func (cfg *Config) FileRootPath() string {
	return "./" + cfg.BuildPath() + "/"
}
