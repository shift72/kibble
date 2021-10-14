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
	"path/filepath"

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

	routeRegistry, err := models.NewRouteRegistryFromConfig(cfg)
	if err != nil {
		log.Errorf("Loading site config failed: %s", err)
		return 1
	}

	renderer := FileRenderer{
		buildPath:  buildPath,
		sourcePath: sourcePath,
	}

	// langRenderer := NewLanguageRenderer(cfg, site)

	// err = langRenderer.PreprocessLanguageFiles(sourcePath)
	// if err != nil {
	// 	return 1
	// }

	renderer.Initialise()

	initSW.Completed()
	errCount := 0
	renderSW := utils.NewStopwatchLevel("render", logging.NOTICE)

	var defaultLanguage string
	var languageConfigs map[string]models.LanguageConfig
	if site.Toggles["translations_api"] {
		defaultLanguage = site.DefaultLanguage
		languageConfigs = site.LanguagesToLanguageConfigs()
	} else {
		defaultLanguage = cfg.DefaultLanguage
		languageConfigs = cfg.Languages
	}

	//changed range from cfg.Langagues - > site.Languages
	for languageObjKey, languageObj := range languageConfigs {

		code := languageObj.Code

		ctx := models.RenderContext{
			RoutePrefix: "",
			Site:        site,
			Language:    createLanguage(defaultLanguage, languageObjKey, code),
		}

		if languageObjKey != defaultLanguage {
			ctx.RoutePrefix = fmt.Sprintf("/%s", languageObjKey)
			i18n.LoadTranslationFile(filepath.Join(sourcePath, ctx.Language.DefinitionFilePath))
		} else {
			i18n.MustLoadTranslationFile(filepath.Join(sourcePath, ctx.Language.DefinitionFilePath))
		}

		renderLangSW := utils.NewStopwatchf("  render language: %s", languageObjKey)

		T, err := i18n.Tfunc(languageObj.Code, defaultLanguage)
		if err != nil {
			log.Errorf("Translation failed: %s", err)
			errCount++
		}

		// set the template view
		renderer.view = models.CreateTemplateView(routeRegistry, T, &ctx, sourcePath)

		for _, route := range routeRegistry.GetAll() {
			renderRouteSW := utils.NewStopwatchf("    render route %s", route.Name)

			// set the route on the render context for datasources
			ctx.Route = route
			if route.ResolvedDataSource != nil {
				errCount += route.ResolvedDataSource.Iterator(ctx, renderer)
			}
			renderRouteSW.Completed()
		}

		renderLangSW.Completed()
	}

	renderSW.Completed()

	return errCount

}

func createLanguage(defaultLanguage string, languageObjKey string, code string) *models.Language {
	return &models.Language{
		Code:               languageObjKey,
		Locale:             code,
		IsDefault:          (lang == defaultLanguage),
		DefinitionFilePath: fmt.Sprintf("%s.all.json", locale),
	}
}
