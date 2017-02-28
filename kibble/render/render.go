package render

import (
	"fmt"
	"time"

	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/datastore"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/nicksnyder/go-i18n/i18n"
)

// Render - render the files
func Render() {

	datastore.Init()

	cfg := config.LoadConfig()

	i18n.MustLoadTranslationFile(fmt.Sprintf("%s.all.json", cfg.Languages[cfg.DefaultLanguage]))

	routeRegistry := models.NewRouteRegistryFromConfig(cfg)

	start := time.Now()
	for lang, locale := range cfg.Languages {
		i18n.LoadTranslationFile(fmt.Sprintf("%s.all.json", locale))

		T, err := i18n.Tfunc(locale, cfg.DefaultLanguage)
		if err != nil {
			fmt.Println(err)
		}

		ctx := models.RenderContext{
			RoutePrefix: "",
		}

		renderer := ConsoleRenderer{
			view: models.CreateTemplateView(&routeRegistry, T, ctx),
		}

		if lang != cfg.DefaultLanguage {
			ctx.RoutePrefix = fmt.Sprintf("/%s", lang)
		}

		for _, route := range routeRegistry.GetAll() {

			ctx.Route = route

			if route.ResolvedDataSouce != nil {
				route.ResolvedDataSouce.Iterator(ctx, renderer)
			}
		}
	}

	stop := time.Now()
	fmt.Printf("Render completed: %s", stop.Sub(start))
}
