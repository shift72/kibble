package models

import "github.com/CloudyKit/jet"

// Film - represents a film
type Film struct {
	ID       int
	Slug     string
	Title    string
	Synopsis string
}

// Renderer - rendering implementation
type Renderer interface {
	Render(route *Route, filePath string, data jet.VarMap)
}

func CreateTemplateView(routeRegistry *RouteRegistry) *jet.Set {
	view := jet.NewHTMLSet("./templates")
	view.AddGlobal("version", "v1.1.145")
	view.AddGlobal("routeTo", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(entity, "")
	})
	view.AddGlobal("routeToWithName", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(entity, routeName)
	})
	view.AddGlobal("routeToSlug", func(slug string) string {
		return routeRegistry.GetRouteForSlug(slug, "")
	})
	view.AddGlobal("routeToSlugWithName", func(slug string, routeName string) string {
		return routeRegistry.GetRouteForSlug(slug, routeName)
	})

	return view
}
