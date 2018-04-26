package models

import (
	"fmt"
	"os"
	"reflect"
)

// Route - represents a route for rendering and
type Route struct {
	Name               string       `json:"name"`
	URLPath            string       `json:"urlPath"`
	TemplatePath       string       `json:"templatePath"`
	DataSource         string       `json:"datasource"`
	ResolvedDataSouce  DataSource   `json:"-"`
	ResolvedEntityType reflect.Type `json:"-"`
	PageSize           int          `json:"pageSize"`
	Pagination         Pagination   `json:"-"`
}

// Clone - create a copy of the route
func (r *Route) Clone() *Route {
	return &Route{
		Name:               r.Name,
		URLPath:            r.URLPath,
		TemplatePath:       r.TemplatePath,
		DataSource:         r.DataSource,
		ResolvedDataSouce:  r.ResolvedDataSouce,
		ResolvedEntityType: r.ResolvedEntityType,
		PageSize:           r.PageSize,
		Pagination:         r.Pagination,
	}
}

// RouteRegistry - stores a list of routes
type RouteRegistry struct {
	routes []*Route
}

// NewRouteRegistry - create a new route registry
func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		routes: make([]*Route, 0),
	}
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
		return ctx.Route.ResolvedDataSouce.GetRouteForEntity(ctx, entity)
	}

	return fmt.Sprintf("!Error. Route not found for entity:%s and route name '%s'", reflect.TypeOf(entity).Name(), routeName)
}

// GetRouteForSlug - finds the route by the name and type and creates a route from it
func (r *RouteRegistry) GetRouteForSlug(ctx RenderContext, slug string, routeName string) string {

	ctx.Route = r.FindBySlugAndRouteName(slug, routeName)

	if ctx.Route != nil {
		return ctx.Route.ResolvedDataSouce.GetRouteForSlug(ctx, slug)
	}

	return fmt.Sprintf("!Error. Route not found for slug:%s and route name '%s'", slug, routeName)
}

// NewRouteRegistryFromConfig - create a new route registry from the config
func NewRouteRegistryFromConfig(config *Config) *RouteRegistry {
	routeRegistry := NewRouteRegistry()

	routeRegistry.routes = make([]*Route, len(config.Routes))

	errorsFound := false

	for i := 0; i < len(config.Routes); i++ {
		route := config.Routes[i]

		route.ResolvedDataSouce = FindDataSource(route.DataSource)
		if route.ResolvedDataSouce != nil {
			route.ResolvedEntityType = route.ResolvedDataSouce.GetEntityType()
		} else {
			fmt.Printf("Unable to find the datasource %s. Check routes registered in site.json\n", route.DataSource)
			errorsFound = true

		}
		routeRegistry.routes[i] = &route
	}

	if errorsFound {
		os.Exit(1)
	}

	return routeRegistry
}
