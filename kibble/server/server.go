package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/datastore"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

var view *jet.Set

// StartNew - start a new server
func StartNew(port int32) {

	view = jet.NewHTMLSet("./templates")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(10 * time.Second))

	loadRoutes(r)

	fmt.Printf("listening on %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func loadRoutes(r chi.Router) {

	var View = jet.NewHTMLSet("./templates")
	View.SetDevelopmentMode(true)

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("kibble online\r\n"))
	})

	for _, route := range *models.AllRoutes {
		switch route.DataSource {
		case "Film":
			r.Get(route.URLPath, routeFilm(route.TemplatePath))
		case "AllFilms":
			r.Get(route.URLPath, routeAllFilms(route.TemplatePath))
			//TODO: add support for tv / season / pages
		default:
			log.Printf("Unknown data source %s\n", route.DataSource)
		}
	}
}

func processRoute(
	w http.ResponseWriter,
	r *http.Request,
	templatePath string,
	dataRequest func(*http.Request) (jet.VarMap, error)) {

	data, err := dataRequest(r)
	if err != nil || data == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, http.StatusText(http.StatusNotFound))
		return
	}

	t, err := view.GetTemplate(templatePath)
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

func routeAllFilms(templateName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		processRoute(w, r, templateName,
			func(r *http.Request) (jet.VarMap, error) {

				f, err := datastore.GetAllFilms()
				if err != nil || f == nil {
					return nil, err
				}

				vars := make(jet.VarMap)
				vars.Set("films", f)
				return vars, nil
			},
		)
	}
}

func routeFilm(templateName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		processRoute(w, r, "film/item.jet",
			func(r *http.Request) (jet.VarMap, error) {

				filmID, err := strconv.Atoi(chi.URLParam(r, "filmID"))
				if err != nil {
					return nil, err
				}

				f, err := datastore.FindByID(filmID)
				if err != nil || f == nil {
					return nil, err
				}

				vars := make(jet.VarMap)
				vars.Set("film", f)
				return vars, nil
			},
		)
	}
}
