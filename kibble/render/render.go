package render

import (
	"bytes"
	"fmt"
	"log"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/datastore"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

var view *jet.Set

// Render - render the files
func Render() {

	//TODO: for the defined languages
	fmt.Println("render called")

	view = jet.NewHTMLSet("./templates")

	for _, route := range *models.AllRoutes {

		fmt.Printf("render route: %s\n", route.URLPath)

		ds := datastore.FindDataSource(route.DataSource)
		if ds == nil {
			log.Printf("Unknown data source %s\n", route.DataSource)
		}

		if ds != nil {
			ds.Iterator(&route, renderToString)
		}
	}

	//TODO: add a stop watch
}

//TODO: make this an interface so we can swap it out
func renderToString(r *models.Route, path string, data jet.VarMap) {

	w := bytes.NewBufferString("")
	w.Write([]byte("--------------------\n"))

	w.Write([]byte(path))

	w.Write([]byte("--------------------\n"))

	t, err := view.GetTemplate(r.TemplatePath)
	if err != nil {
		w.Write([]byte("Template error\n"))
		w.Write([]byte(err.Error()))
		return
	}

	if err = t.Execute(w, data, nil); err != nil {
		w.Write([]byte("Execute error\n"))
		w.Write([]byte(err.Error()))
	}

	fmt.Println(w)
}
