package models

import (
	"reflect"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/version"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/russross/blackfriday"
)

// CreateTemplateView - create a template view
func CreateTemplateView(routeRegistry *RouteRegistry, trans i18n.TranslateFunc, ctx RenderContext, templatePath string) *jet.Set {
	view := jet.NewHTMLSet(templatePath)
	view.AddGlobal("version", version.Version)
	view.AddGlobal("routeTo", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(ctx, entity, "")
	})
	view.AddGlobal("routeToWithName", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(ctx, entity, routeName)
	})
	view.AddGlobal("routeToSlug", func(slug string) string {
		return routeRegistry.GetRouteForSlug(ctx, slug, "")
	})
	view.AddGlobal("canonicalRouteToSlug", func(slug string) string {
		// override the route prefix
		ctxClone := ctx
		ctxClone.RoutePrefix = ""
		return routeRegistry.GetRouteForSlug(ctxClone, slug, "")
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

			log.Errorf("WARN: translating %s found unrecognised type %s", translationID, reflect.TypeOf(args[0]))
		}
		return trans(translationID)

	})

	view.AddGlobal("config", func(key string) string {
		return ctx.Site.Config[key]
	})

	view.AddGlobal("isEnabled", func(key string) bool {
		return ctx.Site.Toggles[key]
	})

	return view
}

// ApplyContentTransforms - add the markdown / sanitization / shortcodes
func ApplyContentTransforms(data string) string {

	//TODO: apply shortcodes

	// apply mark down
	unsafe := blackfriday.MarkdownCommon([]byte(data))

	// return string(unsafe)
	// apply sanitization
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return string(html)
}
