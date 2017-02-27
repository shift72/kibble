package models

import (
	"fmt"
	"reflect"
)

// Route - represents a route for rendering and
type Route struct {
	Name               string
	URLPath            string
	TemplatePath       string
	DataSource         string
	ResolvedDataSouce  DataSource
	ResolvedEntityType reflect.Type
}

// RouteRegistry - stores a list of routes
type RouteRegistry struct {
	routes []*Route
}

// NewRouteRegistry - create a new route registry
func NewRouteRegistry() RouteRegistry {
	return RouteRegistry{
		routes: make([]*Route, 0, 10),
	}
}

// Add - add a route to route registry
func (r *RouteRegistry) Add(route *Route) {

	// find data source
	route.ResolvedDataSouce = FindDataSource(route.DataSource)

	if route.ResolvedDataSouce == nil {
		fmt.Printf("Unable to find the datasource %s\n", route.DataSource)
		return
	}

	route.ResolvedEntityType = route.ResolvedDataSouce.GetEntityType()

	r.routes = append(r.routes, route)
}

// FindByName - find the route by the name
func (r *RouteRegistry) FindByName(name string) *Route {
	for _, r := range r.routes {
		if r.Name == name {
			return r
		}
	}
	return nil
}

// FindByTypeAndRouteName - find a route given a type of entity and the route name(optional)
func (r *RouteRegistry) FindByTypeAndRouteName(entityType reflect.Type, routeName string) *Route {
	for _, r := range r.routes {
		if r.ResolvedEntityType == entityType &&
			(routeName == "" || r.Name == routeName) {
			return r
		}
	}
	return nil
}

// FindBySlugAndRouteName - find a route given a slug and the route name(optional)
func (r *RouteRegistry) FindBySlugAndRouteName(slug string, routeName string) *Route {
	for _, r := range r.routes {
		if r.ResolvedDataSouce.IsSlugMatch(slug) &&
			(routeName == "" || r.Name == routeName) {
			return r
		}
	}
	return nil
}

// GetAll - return all routes
func (r *RouteRegistry) GetAll() []*Route {
	return r.routes
}

// GetRouteForEntity - finds the route by the name and type and creates a route from it
func (r *RouteRegistry) GetRouteForEntity(ctx RenderContext, entity interface{}, routeName string) string {

	ctx.Route = r.FindByTypeAndRouteName(reflect.TypeOf(entity), routeName)

	if ctx.Route != nil {
		// fmt.Printf("Found route, name:%s, path: %s\n", foundRoute.Name, foundRoute.URLPath)
		return ctx.Route.ResolvedDataSouce.GetRouteForEntity(ctx, entity)
	}

	return fmt.Sprintf("!Error. Route not found for entity:%s and route name %v", reflect.TypeOf(entity).Name(), routeName)
}

// GetRouteForSlug - finds the route by the name and type and creates a route from it
func (r *RouteRegistry) GetRouteForSlug(ctx RenderContext, slug string, routeName string) string {

	ctx.Route = r.FindBySlugAndRouteName(slug, routeName)

	if ctx.Route != nil {
		// fmt.Printf("Found route, name:%s, path: %s\n", foundRoute.Name, foundRoute.URLPath)
		return ctx.Route.ResolvedDataSouce.GetRouteForSlug(ctx, slug)
	}

	return fmt.Sprintf("!Error. Route not found for slug:%s and route name %v", slug, routeName)
}

// NewRouteRegistryFromConfig - create a new route registry from the config
func NewRouteRegistryFromConfig(config *Config) RouteRegistry {
	routeRegistry := NewRouteRegistry()

	for _, r := range config.Routes {
		routeRegistry.Add(&r)
	}
	return routeRegistry
}

// DefaultRouteRegistry - to be replaced with one that loads routes from a file
// func DefaultRouteRegistry() RouteRegistry {
//
// 	routeRegistry := NewRouteRegistry()
// 	routeRegistry.Add(&Route{
// 		Name:         "filmIndex",
// 		URLPath:      "/film",
// 		TemplatePath: "film/index.jet",
// 		DataSource:   "FilmCollection",
// 	})
//
// 	routeRegistry.Add(&Route{
// 		Name:         "filmItem",
// 		URLPath:      "/film/:filmID",
// 		TemplatePath: "film/item.jet",
// 		DataSource:   "Film",
// 	})
//
// 	routeRegistry.Add(&Route{
// 		Name:         "filmItemPartial",
// 		URLPath:      "/film/:filmID/partial.html",
// 		TemplatePath: "film/partial.jet",
// 		DataSource:   "Film",
// 	})
// 	return routeRegistry
// }
