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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	err = preprocessLanguageFiles(site, cfg, sourcePath)
	if err != nil {
		return 1
	}

	renderer.Initialise()

	initSW.Completed()

	errCount := 0

	renderSW := utils.NewStopwatchLevel("render", logging.NOTICE)

	for languageObjKey, languageObj := range cfg.Languages {

		code := languageObj.Code

		ctx := models.RenderContext{
			RoutePrefix: "",
			Site:        site,
			Language:    formatContextLanguage(site.Toggles["translations_api"], cfg, languageObjKey, code),
		}

		if languageObjKey != cfg.DefaultLanguage {
			ctx.RoutePrefix = fmt.Sprintf("/%s", languageObjKey)
			i18n.LoadTranslationFile(filepath.Join(sourcePath, ctx.Language.DefinitionFilePath))
		} else {
			i18n.MustLoadTranslationFile(filepath.Join(sourcePath, ctx.Language.DefinitionFilePath))
		}

		renderLangSW := utils.NewStopwatchf("  render language: %s", languageObjKey)

		T, err := i18n.Tfunc(languageObj.Code, cfg.DefaultLanguage)
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

func createLanguage(cfg *models.Config, languageObjKey string, code string) *models.Language {
	return &models.Language{
		Code:               languageObjKey,
		Locale:             code,
		IsDefault:          (languageObjKey == cfg.DefaultLanguage),
		DefinitionFilePath: formatLanguageFilename(code),
	}
}

func overrideLanguages(site *models.Site, cfg *models.Config) error {
	defaultLanguageOverride := func(isDefault bool, langCode string) string {
		if isDefault {
			return ""
		}
		return langCode
	}

	if site.Toggles["translations_api"] {

		site.Languages = make([]models.Language, 0)

		cfg.DefaultLanguage = ""
		cfg.Languages = make(map[string]models.LanguageConfig)

		languages, err := api.LoadAllLanguages(cfg)
		if err != nil {
			log.Errorf("Failed to get languages: %s", err)
			return err
		}

		cfg.DefaultLanguage = formatPathLocale(languages.DefaultLanguage.Code)

		for _, lang := range languages.SupportedLanguages {
			langCode := formatPathLocale(lang.Code)
			langLabel := formatPathLocale(lang.Label)

			isDefault := langCode == cfg.DefaultLanguage

			site.Languages = append(site.Languages, models.Language{
				IsDefault: isDefault,
				Code:      defaultLanguageOverride(isDefault, langCode),
				Name:      langLabel,
			})

			cfg.Languages[langCode] = models.LanguageConfig{
				Code: langCode,
				Name: langLabel,
			}
		}
	}

	return nil
}

func obtainTranslations(site *models.Site, cfg *models.Config) (api.TranslationsV1, error) {
	if site.Toggles["translations_api"] {

		translations, err := api.LoadAllTranslations(cfg)
		if err != nil {
			return nil, err
		}

		for key := range translations {
			translations[formatPathLocale(key)] = translations[key]
		}

		return translations, nil
	}

	return nil, nil
}

func formatPathLocale(code string) string {
	dashedCode := strings.ReplaceAll(code, "_", "-")
	return strings.ToLower(dashedCode)
}

func writeFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Errorf("%s", err)
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		log.Errorf("%s", err)
		file.Close()
		return err
	}

	return file.Close()
}

func formatLanguageFilename(code string) string {
	return fmt.Sprintf("%s.all.json", code)
}

func preprocessLanguageFiles(site *models.Site, cfg *models.Config, sourcePath string) error {
	err := overrideLanguages(site, cfg)
	if err != nil {
		return err
	}

	apiTranslations, err := obtainTranslations(site, cfg)
	if err != nil {
		log.Errorf("Failed to get translations: %s", err)
		return err
	}

	if apiTranslations != nil {
		for _, languageObj := range cfg.Languages {

			code := languageObj.Code

			filename := formatLanguageFilename(code)

			file, err := json.Marshal(apiTranslations[code])
			if err != nil {
				log.Errorf("Failed to marshal translations json %s: %s", code, err)
				return err
			}

			err = writeFile(filepath.Join(sourcePath, filename), file)
			if err != nil {
				log.Errorf("Failed to write translations files: %s", err)
				return err
			}
		}
	}
	return nil
}

func formatContextLanguage(translationsAPIEnabled bool, cfg *models.Config, languageObjKey string, code string) *models.Language {
	if translationsAPIEnabled {
		return &models.Language{
			Code:               formatPathLocale(languageObjKey),
			Locale:             formatPathLocale(code),
			IsDefault:          (formatPathLocale(languageObjKey) == cfg.DefaultLanguage),
			DefinitionFilePath: formatLanguageFilename(code),
		}
	}
	return createLanguage(cfg, languageObjKey, code)
}
