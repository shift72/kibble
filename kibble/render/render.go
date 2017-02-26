package render

import (
	"fmt"
	"log"
	"time"

	"github.com/indiereign/shift72-kibble/kibble/datastore"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// Render - render the files
func Render() {

	//TODO: for the defined languages
	datastore.Init()

	routeRegistry := models.DefaultRouteRegistry()

	renderer := ConsoleRenderer{
		view: models.CreateTemplateView(&routeRegistry),
	}

	start := time.Now()
	for _, route := range routeRegistry.GetAll() {

		ds := models.FindDataSource(route.DataSource)
		if ds == nil {
			log.Printf("Unknown data source %s\n", route.DataSource)
		}

		if ds != nil {
			ds.Iterator(route, renderer)
		}
	}
	stop := time.Now()
	fmt.Printf("Render completed: %s", stop.Sub(start))
}
