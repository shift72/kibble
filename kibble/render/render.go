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
	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/indiereign/shift72-kibble/kibble/perf"
	"github.com/nicksnyder/go-i18n/i18n"
	logging "github.com/op/go-logging"
)

var rootPath = path.Join(".kibble", "build")
var staticFolder = "static"

// Watch -
func Watch(runAsAdmin bool, port int32) {

	liveReload := LiveReload{}
	liveReload.StartLiveReload(port, func() {
		// re-render
		Render(runAsAdmin)
	})

	cfg := config.LoadConfig(runAsAdmin)
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
func Render(runAsAdmin bool) {

	initSW := perf.NewStopwatch("load")

	cfg := config.LoadConfig(runAsAdmin)

	api.CheckAdminCredentials(cfg, runAsAdmin)

	site, err := api.LoadSite(cfg)
	if err != nil {
		log.Errorf("Site load failed: %s", err)
		return
	}

	routeRegistry := models.NewRouteRegistryFromConfig(cfg)

	renderer := FileRenderer{
		rootPath: rootPath,
	}

	renderer.Initialise()

	initSW.Completed()

	sassSW := perf.NewStopwatch("sass")
	err = Sass(
		path.Join("styles", "main.scss"),
		path.Join(rootPath, "styles", "main.css"))
	if err != nil {
		log.Errorf("Sass rendering failed: %s", err)
		return
	}
	sassSW.Completed()

	renderSW := perf.NewStopwatchLevel("render", logging.NOTICE)
	for lang, locale := range cfg.Languages {

		renderLangSW := perf.NewStopwatchf("  render language: %s", lang)
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

		renderFilesSW := perf.NewStopwatch("  render files")
		for _, f := range files {
			filePath := path.Join(ctx.RoutePrefix, strings.Replace(f, ".jet", "", 1))

			route := &models.Route{
				TemplatePath: f,
			}

			data := jet.VarMap{}
			data.Set("site", site)
			renderer.Render(route, filePath, data)
		}
		renderFilesSW.Completed()

		for _, route := range routeRegistry.GetAll() {
			renderRouteSW := perf.NewStopwatchf("    render route %s", route.Name)
			ctx.Route = route
			if route.ResolvedDataSouce != nil {
				route.ResolvedDataSouce.Iterator(ctx, renderer)
			}
			renderRouteSW.Completed()
		}

		renderLangSW.Completed()
	}

	renderSW.Completed()
}
