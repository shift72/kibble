package render

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/api"
	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/datastore"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/nicksnyder/go-i18n/i18n"
)

var rootPath = path.Join(".", ".kibble", "build")
var publicFolder = "public"

// Watch -
func Watch(runAsAdmin bool, verbose bool, port int32) {

	liveReload := LiveReload{}
	liveReload.StartLiveReload(port, func() {
		// re-render
		Render(runAsAdmin, verbose)
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
		fmt.Println(err)
		os.Exit(1)
	}
}

// Render - render the files
func Render(runAsAdmin bool, verbose bool) {

	datastore.Init()

	cfg := config.LoadConfig(runAsAdmin)

	api.CheckAdminCredentials(cfg, runAsAdmin)

	site, err := api.LoadSite(cfg)
	if err != nil {
		fmt.Printf("Site load failed: %s", err)
		return
	}

	routeRegistry := models.NewRouteRegistryFromConfig(cfg)

	renderer := FileRenderer{
		rootPath:    rootPath,
		showSummary: verbose,
	}
	renderer.Initialise()

	start := time.Now()

	err = Sass(
		path.Join("styles", "main.scss"),
		path.Join(rootPath, "styles", "main.css"))
	if err != nil {
		fmt.Printf("Sass rendering failed: %s", err)
		return
	}

	fmt.Printf("Sass render time: %s\n", time.Now().Sub(start))

	for lang, locale := range cfg.Languages {

		T, err := i18n.Tfunc(locale, cfg.DefaultLanguage)
		if err != nil {
			fmt.Println(err)
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
		for _, f := range files {
			filePath := path.Join(ctx.RoutePrefix, strings.Replace(f, ".jet", "", 1))

			route := &models.Route{
				TemplatePath: f,
			}

			data := jet.VarMap{}
			data.Set("site", site)
			renderer.Render(route, filePath, data)
		}

		for _, route := range routeRegistry.GetAll() {

			ctx.Route = route
			if route.ResolvedDataSouce != nil {
				route.ResolvedDataSouce.Iterator(ctx, renderer)
			}
		}
	}

	stop := time.Now()

	fmt.Printf("\nTotal render time: %s\n", stop.Sub(start))
}
