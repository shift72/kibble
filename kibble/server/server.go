package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/api"
	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/datastore"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

var site *models.Site

// StartNew - start a new server
func StartNew(port int32) {

	datastore.Init()

	cfg := config.LoadConfig()

	var err error
	site, err = api.LoadSite(cfg)
	if err != nil {
		fmt.Printf("Site load failed: %s", err)
		return
	}

	routeRegistry := models.NewRouteRegistryFromConfig(cfg)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)
	r.Use(FileMiddleware(cfg, site, routeRegistry))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(10 * time.Second))

	loadRoutes(r, routeRegistry, cfg)

	fmt.Printf("listening on %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func loadRoutes(r chi.Router, routeRegistry *models.RouteRegistry, cfg *models.Config) {

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("kibble online\r\n"))
	})

	//TODO: sort the routes, put the collections at the end
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

		ctx := models.RenderContext{
			Route:       route,
			RoutePrefix: "",
			Site:        site,
			Language:    chi.URLParam(req, "lang"),
		}

		data, err := route.ResolvedDataSouce.Query(ctx, req)
		if err != nil || data == nil {
			render.Status(req, http.StatusNotFound)
			render.JSON(w, req, http.StatusText(http.StatusNotFound))
			return
		}

		renderContext(cfg, routeRegistry, ctx, "./templates", data, w, req)
	}
}

// FileMiddleware -
func FileMiddleware(cfg *models.Config, site *models.Site, routeRegistry *models.RouteRegistry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			pwd, _ := os.Getwd()

			//TODO: check the languages
			path := fmt.Sprintf("%s%s.jet", pwd, r.RequestURI)
			templatePath := fmt.Sprintf("%s.jet", r.RequestURI)

			// check if the request + jet file exists
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("not exists", err)
				next.ServeHTTP(w, r)
				return
			}

			ctx := models.RenderContext{
				Route: &models.Route{
					TemplatePath: templatePath,
				},
				RoutePrefix: "",
				Site:        site,
				Language:    cfg.DefaultLanguage,
			}

			data := jet.VarMap{}
			data.Set("site", site)

			renderContext(cfg, routeRegistry, ctx, "./", data, w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func renderContext(
	cfg *models.Config,
	routeRegistry *models.RouteRegistry,
	ctx models.RenderContext,
	templatePath string,
	data jet.VarMap,
	w http.ResponseWriter,
	req *http.Request) {

	// fmt.Printf("lang:%s\n", lang)
	// fmt.Printf("locale:%s\n", cfg.Languages[lang])
	// fmt.Printf("locale default:%s\n", cfg.Languages[cfg.DefaultLanguage])

	T, err := i18n.Tfunc(cfg.Languages[ctx.Language], cfg.Languages[cfg.DefaultLanguage])
	if err != nil {
		fmt.Println(err)
	}

	if ctx.Language != "" && ctx.Language != cfg.DefaultLanguage {
		ctx.RoutePrefix = fmt.Sprintf("/%s", ctx.Language)
	}

	view := models.CreateTemplateView(routeRegistry, T, ctx, templatePath)
	view.SetDevelopmentMode(true)

	t, err := view.GetTemplate(ctx.Route.TemplatePath)
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
