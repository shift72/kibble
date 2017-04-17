package render

import (
	"fmt"
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
		rootPath:    "./.kibble/build",
		showSummary: verbose,
	}
	renderer.Initialise()

	start := time.Now()
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

		// set the template view
		renderer.view = models.CreateTemplateView(routeRegistry, T, ctx, "./")

		if lang != cfg.DefaultLanguage {
			ctx.RoutePrefix = fmt.Sprintf("/%s", lang)
		}

		// render static files
		files, _ := filepath.Glob("*.jet")
		for _, f := range files {

			filePath := fmt.Sprintf("%s/%s", ctx.RoutePrefix, strings.Replace(f, ".jet", "", 1))

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

	fmt.Printf("\nRendered: %s", stop.Sub(start))
}
