package server

import (
	"fmt"
	"log"
	"net/http"
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

	view.SetDevelopmentMode(true)

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("kibble online\r\n"))
	})

	for _, route := range *models.AllRoutes {

		ds := datastore.FindDataSource(route.DataSource)
		if ds == nil {
			log.Printf("Unknown data source %s\n", route.DataSource)
		}

		if ds != nil {
			r.Get(route.URLPath, routeToDataSoure(route.TemplatePath, ds))
		}
	}
}

func routeToDataSoure(templateName string, ds *datastore.DataSource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		processRoute(w, r, templateName, ds)
	}
}

func processRoute(
	w http.ResponseWriter,
	req *http.Request,
	templatePath string,
	ds *datastore.DataSource) {

	data, err := ds.Query(req)
	if err != nil || data == nil {
		render.Status(req, http.StatusNotFound)
		render.JSON(w, req, http.StatusText(http.StatusNotFound))
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
