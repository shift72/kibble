package models

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

// AllRoutes - create the routes
var AllRoutes = &[]Route{
	{
		URLPath:      "/film",
		TemplatePath: "film/index.jet",
		DataSource:   "AllFilms",
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
