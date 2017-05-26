package render

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/api"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/utils"
	"github.com/nicksnyder/go-i18n/i18n"
	logging "github.com/op/go-logging"
)

var staticFolder = "static"

// Watch -
func Watch(rootPath string, cfg *models.Config, port int32, logReader utils.LogReader) {

	liveReload := LiveReload{logReader: logReader}
	liveReload.StartLiveReload(port, func() {
		// re-render
		logReader.Clear()
		Render(rootPath, cfg)
	})

	proxy := NewProxy(cfg.SiteURL)

	// server
	mux := http.NewServeMux()
	mux.HandleFunc("/kibble/live_reload", liveReload.Handler)
	mux.Handle("/",
		proxy.GetMiddleware(
			liveReload.GetMiddleware(
				http.FileServer(
					http.Dir(rootPath)))))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Errorf("Web server failed: %s", err)
		os.Exit(1)
	}
}

// Render - render the files
func Render(rootPath string, cfg *models.Config) error {

	initSW := utils.NewStopwatch("load")

	api.CheckAdminCredentials(cfg)

	site, err := api.LoadSite(cfg)
	if err != nil {
		return err
	}

	routeRegistry := models.NewRouteRegistryFromConfig(cfg)

	renderer := FileRenderer{
		rootPath: rootPath,
	}

	renderer.Initialise()

	initSW.Completed()

	errCount := 0

	renderSW := utils.NewStopwatchLevel("render", logging.NOTICE)
	for lang, locale := range cfg.Languages {

		renderLangSW := utils.NewStopwatchf("  render language: %s", lang)
		T, err := i18n.Tfunc(locale, cfg.DefaultLanguage)
		if err != nil {
			log.Errorf("Translation failed: %s", err)
		}

		ctx := models.RenderContext{
			RoutePrefix: "",
			Site:        site,
			Language:    lang,
		}

		if lang != cfg.DefaultLanguage {
			ctx.RoutePrefix = fmt.Sprintf("/%s", lang)
		}

		// set the template view
		renderer.view = models.CreateTemplateView(routeRegistry, T, ctx, "./")

		// render static files
		files, _ := filepath.Glob("*.jet")

		renderFilesSW := utils.NewStopwatch("  render files")
		for _, f := range files {
			filePath := path.Join(ctx.RoutePrefix, strings.Replace(f, ".jet", "", 1))

			route := &models.Route{
				TemplatePath: f,
			}

			data := jet.VarMap{}
			data.Set("site", site)
			errCount += renderer.Render(route, filePath, data)
		}
		renderFilesSW.Completed()

		for _, route := range routeRegistry.GetAll() {
			renderRouteSW := utils.NewStopwatchf("    render route %s", route.Name)
			ctx.Route = route
			if route.ResolvedDataSouce != nil {
				route.ResolvedDataSouce.Iterator(ctx, renderer)
			}
			renderRouteSW.Completed()
		}

		renderLangSW.Completed()
	}

	renderSW.Completed()

	log.Debug("error count %d", errCount)
	return nil
}
