//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package models

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/version"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/russross/blackfriday.v2"
)

var templateTagRegex = regexp.MustCompile("(?U:{{.+}})+")
var shortCodeView *jet.Set

// CreateTemplateView - create a template view
func CreateTemplateView(routeRegistry *RouteRegistry, trans i18n.TranslateFunc, ctx *RenderContext, templatePath string) *jet.Set {
	view := jet.NewHTMLSet(templatePath)
	view.AddGlobal("version", version.Version)
	view.AddGlobal("lang", ctx.Language)
	view.AddGlobal("routeTo", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(*ctx, entity, "")
	})
	view.AddGlobal("routeToWithName", func(entity interface{}, routeName string) string {
		return routeRegistry.GetRouteForEntity(*ctx, entity, routeName)
	})
	view.AddGlobal("routeToSlug", func(slug string) string {
		return routeRegistry.GetRouteForSlug(*ctx, slug, "")
	})
	view.AddGlobal("canonicalRouteToSlug", func(slug string) string {
		// override the route prefix
		ctxClone := *ctx
		ctxClone.RoutePrefix = ""
		return routeRegistry.GetRouteForSlug(ctxClone, slug, "")
	})
	view.AddGlobal("routeToSlugWithName", func(slug string, routeName string) string {
		return routeRegistry.GetRouteForSlug(*ctx, slug, routeName)
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

			argType := reflect.TypeOf(args[0])
			argTypeName := argType.String()

			if argTypeName == "string" {
				if f, ok := args[0].(string); ok && len(f) == 0 {
					log.Warningf("WARN: found empty string parameter while translating \"%s\". Set a default or ensure the argument is not empty in %s", translationID, ctx.Route.TemplatePath)
				}
			}

			if argTypeName == "string" || strings.Contains(argTypeName, "int") || argTypeName == "map[string]interface {}" {
				return trans(translationID, args[0])
			}

			if f, ok := args[0].(float64); ok {
				return trans(translationID, int(f))
			}

			log.Errorf("WARN: translating %s found unrecognised type %s", translationID, argType)
		}
		return trans(translationID)
	})

	view.AddGlobal("config", func(key string, args ...string) string {
		if s, ok := ctx.Site.Config[key]; ok {
			return s
		}

		if len(args) > 0 {
			return args[0]
		}

		return ""
	})

	view.AddGlobal("configInt", func(key string, args ...int) int {
		return ctx.Site.Config.GetInt(key, args...)
	})

	view.AddGlobal("isEnabled", func(key string) bool {
		return ctx.Site.Toggles[key]
	})

	return view
}

// ApplyContentTransforms - add the markdown / sanitization / shortcodes
func ApplyContentTransforms(data string) string {

	// apply mark down
	unsafe := blackfriday.Run([]byte(data))

	// apply the templates
	return insertTemplates(string(unsafe))
}

// insertTemplates applys any templates and sanitise the output
func insertTemplates(data string) string {
	var p string

	matches := templateTagRegex.FindAllStringSubmatchIndex(data, -1)

	cleaner := bluemonday.UGCPolicy()
	cleaner.AddTargetBlankToFullyQualifiedLinks(true)
	cleaner.AllowAttrs("target").OnElements("a")

	c := len(matches)
	if c > 0 {
		p = ""
		for i := 0; i < c; i++ {
			if i == 0 {
				p = p + cleaner.Sanitize(data[:matches[i][0]]) +
					processTemplateTag(data[matches[i][0]:matches[i][1]])
			}
			if i == c-1 {
				p = p + cleaner.Sanitize(data[matches[i][1]:])
			} else {
				p = p + cleaner.Sanitize(data[matches[i][1]:matches[i+1][0]]) +
					processTemplateTag(data[matches[i+1][0]:matches[i+1][1]])
			}
		}
	} else {
		p = cleaner.Sanitize(data)
	}

	return p
}

// ConfigureShortcodeTemplatePath sets the directory where the short codes
// will be loaded from
func ConfigureShortcodeTemplatePath(cfg *Config) {

	if shortCodeView == nil {
		// get the template view
		shortCodeView = jet.NewHTMLSet(cfg.ShortCodePath())

		// built-in templates
		shortCodeView.LoadTemplate("echo.jet", "<div class=\"echo\">slug:{{slug}}</div>")
		shortCodeView.LoadTemplate("youtube.jet", `
<div {{isset(class) ? "class=\"" + class + "\"" : "style=\"position: relative; padding-bottom: 56.25%; padding-top: 30px; height: 0; overflow: hidden;\"" | raw }} >
<iframe src="//www.youtube.com/embed/{{id}}" {{isset(class) ? "class=\"" + class + "\"" : "style=\"position: absolute; top: 0; left: 0; width: 100%; height: 100%;\"" | raw }}{{if isset(autoplay) && autoplay=="true" }} autoplay=1{{end}} allowfullscreen frameborder="0"></iframe>
</div>`)
	}
}

func processTemplateTag(templateTag string) string {

	templateName, data, err := parseParameters(templateTag)

	w := bytes.NewBufferString("")
	templatePath := fmt.Sprintf("%s.jet", templateName)
	t, err := shortCodeView.GetTemplate(templatePath)
	if err != nil {
		log.Error("Template load error. Loading %s %s", templatePath, err)
		return "Err"
	}

	if err = t.Execute(w, *data, nil); err != nil {
		w.WriteString("<pre>")
		w.WriteString(err.Error())
		w.WriteString("</pre>")
		log.Errorf("Template execute error: %s", err)
	}

	return string(w.Bytes())
}

// parse the template tag into the template and arguments
func parseParameters(templateTag string) (string, *jet.VarMap, error) {

	parameters := strings.Fields(strings.Trim(templateTag, "{} "))

	data := jet.VarMap{}

	if len(parameters) > 1 {
		for _, p := range parameters[1:] {
			s := strings.Split(p, "=")
			if len(s) == 2 {
				data.Set(s[0], s[1])
			}
		}
	}

	return parameters[0], &data, nil
}
