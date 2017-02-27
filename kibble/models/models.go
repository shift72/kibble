package models

import (
	"github.com/CloudyKit/jet"
	"github.com/nicksnyder/go-i18n/i18n"
)

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

// CreateTemplateView - create a template view
func CreateTemplateView(routeRegistry *RouteRegistry, trans i18n.TranslateFunc, ctx RenderContext) *jet.Set {
	view := jet.NewHTMLSet("./templates")
	view.AddGlobal("version", "v1.1.145")
	view.AddGlobal("routeTo", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(ctx, entity, "")
	})
	view.AddGlobal("routeToWithName", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(ctx, entity, routeName)
	})
	view.AddGlobal("routeToSlug", func(slug string) string {
		return routeRegistry.GetRouteForSlug(ctx, slug, "")
	})
	view.AddGlobal("routeToSlugWithName", func(slug string, routeName string) string {
		return routeRegistry.GetRouteForSlug(ctx, slug, routeName)
	})
	view.AddGlobal("i18n", trans)

	return view
}
