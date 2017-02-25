package models

import "github.com/CloudyKit/jet"

// Film - represents a film
type Film struct {
	ID       int
	Slug     string
	Title    string
	Synopsis string
}

// Route - represents a route for rendering and
type Route struct {
	URLPath      string
	TemplatePath string
	DataSource   string
}

// Renderer - rendering implementation
type Renderer interface {
	Render(route *Route, filePath string, data jet.VarMap)
}

// AllRoutes - create the routes
var AllRoutes = &[]Route{
	{
		URLPath:      "/film",
		TemplatePath: "film/index.jet",
		DataSource:   "FilmCollection",
	},
	{
		URLPath:      "/film/:filmID",
		TemplatePath: "film/item.jet",
		DataSource:   "Film",
	},
	{
		URLPath:      "/film/:filmID/partial.html",
		TemplatePath: "film/partial.jet",
		DataSource:   "Film",
	},
}
