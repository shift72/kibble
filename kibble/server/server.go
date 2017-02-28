package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/datastore"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

// StartNew - start a new server
func StartNew(port int32) {

	datastore.Init()

	cfg := config.LoadConfig()
	routeRegistry := models.NewRouteRegistryFromConfig(cfg)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(10 * time.Second))

	loadRoutes(r, &routeRegistry, cfg)

	fmt.Printf("listening on %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func loadRoutes(r chi.Router, routeRegistry *models.RouteRegistry, cfg *models.Config) {

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("kibble online\r\n"))
	})

	for _, route := range routeRegistry.GetAll() {
		if route.ResolvedDataSouce != nil {
			r.Get(route.URLPath, routeToDataSoure(route, routeRegistry, cfg))
			r.Get("/:lang"+route.URLPath, routeToDataSoure(route, routeRegistry, cfg))
		} else {
			log.Printf("Route skipped, unknown data source %s\n", route.DataSource)
		}
	}
}

func routeToDataSoure(route *models.Route, routeRegistry *models.RouteRegistry, cfg *models.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		lang := chi.URLParam(req, "lang")

		// fmt.Printf("lang:%s\n", lang)
		// fmt.Printf("locale:%s\n", cfg.Languages[lang])
		// fmt.Printf("locale default:%s\n", cfg.Languages[cfg.DefaultLanguage])

		T, err := i18n.Tfunc(cfg.Languages[lang], cfg.Languages[cfg.DefaultLanguage])
		if err != nil {
			fmt.Println(err)
		}

		ctx := models.RenderContext{
			Route:       route,
			RoutePrefix: fmt.Sprintf("/%s", lang),
		}

		view := models.CreateTemplateView(routeRegistry, T, ctx)
		view.SetDevelopmentMode(true)

		data, err := route.ResolvedDataSouce.Query(req)
		if err != nil || data == nil {
			render.Status(req, http.StatusNotFound)
			render.JSON(w, req, http.StatusText(http.StatusNotFound))
			return
		}

		t, err := view.GetTemplate(route.TemplatePath)
		if err != nil {
			w.Write([]byte("Template error\n"))
			w.Write([]byte(err.Error()))
			return
		}

		if err = t.Execute(w, data, nil); err != nil {
			w.Write([]byte("Execute error\n"))
			w.Write([]byte(err.Error()))
		}
	}
}
