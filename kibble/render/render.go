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

package render

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet"
	"kibble/api"
	"kibble/models"
	"kibble/utils"
	"github.com/nicksnyder/go-i18n/i18n"
	logging "github.com/op/go-logging"
)

var staticFolder = "static"

// Watch -
func Watch(sourcePath string, buildPath string, cfg *models.Config, port int32, logReader utils.LogReader, reloadBrowserOnChange bool) {

	liveReload := LiveReload{
		logReader:                 logReader,
		sourcePath:                sourcePath,
		config:                    cfg.LiveReload,
		reloadBrowserOnFileChange: reloadBrowserOnChange,
	}

	liveReload.StartLiveReload(port, func() {
		// re-render
		logReader.Clear()
		Render(sourcePath, buildPath, cfg)
	})

	proxy := NewProxy(cfg.SiteURL, cfg.ProxyPatterns)

	// server
	mux := http.NewServeMux()
	if reloadBrowserOnChange {
		mux.HandleFunc("/kibble/live_reload", liveReload.Handler)
	}

	mux.Handle("/",
		proxy.GetMiddleware(
			liveReload.GetMiddleware(
				http.FileServer(
					http.Dir(buildPath)))))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Errorf("Web server failed: %s", err)
		os.Exit(1)
	}
}

// Render - render the files
func Render(sourcePath string, buildPath string, cfg *models.Config) int {

	initSW := utils.NewStopwatch("load")

	api.CheckAdminCredentials(cfg)

	models.ConfigureShortcodeTemplatePath(cfg)

	site, err := api.LoadSite(cfg)
	if err != nil {
		log.Errorf("Loading site config failed: %s", err)
		return 1
	}

	routeRegistry := models.NewRouteRegistryFromConfig(cfg)

	renderer := FileRenderer{
		buildPath:  buildPath,
		sourcePath: sourcePath,
	}

	renderer.Initialise()

	initSW.Completed()

	errCount := 0

	renderSW := utils.NewStopwatchLevel("render", logging.NOTICE)
	for lang, localeObj := range cfg.Languages {

		locale := localeObj["code"]

		ctx := models.RenderContext{
			RoutePrefix: "",
			Site:        site,
			Language:    createLanguage(cfg, lang, locale),
		}

		if lang != cfg.DefaultLanguage {
			ctx.RoutePrefix = fmt.Sprintf("/%s", lang)
			i18n.LoadTranslationFile(filepath.Join(sourcePath, ctx.Language.DefinitionFilePath))
		} else {
			i18n.MustLoadTranslationFile(filepath.Join(sourcePath, ctx.Language.DefinitionFilePath))
		}

		renderLangSW := utils.NewStopwatchf("  render language: %s", lang)
		T, err := i18n.Tfunc(locale, cfg.DefaultLanguage)
		if err != nil {
			log.Errorf("Translation failed: %s", err)
			errCount++
		}

		// set the template view
		renderer.view = models.CreateTemplateView(routeRegistry, T, &ctx, sourcePath)

		// render static files
		files, _ := filepath.Glob(filepath.Join(sourcePath, "*.jet"))

		// set the route on the render context for static template
		ctx.Route = &models.Route{
			Name:         "static",
			URLPath:      "",
			TemplatePath: "",
		}

		renderFilesSW := utils.NewStopwatch("  render files")
		for _, f := range files {
			// jet prefers relative template paths, so lets make it relativeish,
			// be removing the `sourcePath` from the start of it.
			relativeFilePath := strings.Replace(f, sourcePath, "", 1)

			// update the route, per file render
			ctx.Route.TemplatePath = relativeFilePath
			ctx.Route.URLPath = path.Join(ctx.RoutePrefix, strings.Replace(relativeFilePath, ".jet", "", 1))

			data := jet.VarMap{}
			data.Set("site", site)
			errCount += renderer.Render(relativeFilePath, ctx.Route.URLPath, data)
		}
		renderFilesSW.Completed()

		for _, route := range routeRegistry.GetAll() {
			renderRouteSW := utils.NewStopwatchf("    render route %s", route.Name)

			// set the route on the render context for datasources
			ctx.Route = route
			if route.ResolvedDataSouce != nil {
				errCount += route.ResolvedDataSouce.Iterator(ctx, renderer)
			}
			renderRouteSW.Completed()
		}

		renderLangSW.Completed()
	}

	renderSW.Completed()

	return errCount
}

func createLanguage(cfg *models.Config, lang string, locale string) *models.Language {
	return &models.Language{
		Code:               lang,
		Locale:             locale,
		IsDefault:          (lang == cfg.DefaultLanguage),
		DefinitionFilePath: fmt.Sprintf("%s.all.json", locale),
	}
}
