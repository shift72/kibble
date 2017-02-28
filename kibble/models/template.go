package models

import (
	"log"
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/nicksnyder/go-i18n/i18n"
)

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
	view.AddGlobal("i18n", func(translationID string, args ...interface{}) string {

		// jet will pass in numeric args as float64
		// need to massage these back into expected types
		/*
					   TranslateFunc returns the translation of the string identified by translationID.

					   If there is no translation for translationID, then the translationID itself is returned. This makes it easy to identify missing translations in your app.

					   If translationID is a non-plural form, then the first variadic argument may be a map[string]interface{} or struct that contains template data.

					   If translationID is a plural form, the function accepts two parameter signatures 1. T(count int, data struct{}) The first variadic argument must be an integer type (int, int8, int16, int32, int64) or a float formatted as a string (e.g. "123.45").
			       The second variadic argument may be a map[string]interface{} or struct{} that contains template data. 2. T(data struct{}) data must be a struct{} or
			       map[string]interface{} that contains a Count field and the template data,
			       Count field must be an integer type (int, int8, int16, int32, int64) or a float formatted as a string (e.g. "123.45").
		*/
		if len(args) == 1 {
			f, ok := args[0].(float64)
			if ok {
				return trans(translationID, int(f))
			}

			s, ok := args[0].(string)
			if ok {
				return trans(translationID, s)
			}

			log.Printf("WARN: translating %s found unrecognised type %s", translationID, reflect.TypeOf(args[0]))
		}
		return trans(translationID)

	})

	return view
}
