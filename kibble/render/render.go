package render

import (
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

	view = jet.NewHTMLSet("./templates")
	view.AddGlobal("version", "v1.1.145")

	renderer := ConsoleRenderer{
		view: view,
	}

	for _, route := range *models.AllRoutes {

		fmt.Printf("render route: %s\n", route.URLPath)

		ds := datastore.FindDataSource(route.DataSource)
		if ds == nil {
			log.Printf("Unknown data source %s\n", route.DataSource)
		}

		if ds != nil {
			ds.Iterator(&route, renderer)
		}
	}

	//TODO: add a stop watch
}
