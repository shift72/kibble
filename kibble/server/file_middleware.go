package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/CloudyKit/jet"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

// FileMiddleware -
func FileMiddleware(cfg *models.Config, site *models.Site, routeRegistry *models.RouteRegistry) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			pwd, _ := os.Getwd()

			//TODO: check the languages
			path := fmt.Sprintf("%s%s.jet", pwd, r.RequestURI)
			templatePath := fmt.Sprintf("%s.jet", r.RequestURI)

			// check if the request + jet file exists
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("not exists", err)
				next.ServeHTTP(w, r)
				return
			}

			ctx := models.RenderContext{
				Route: &models.Route{
					TemplatePath: templatePath,
				},
				RoutePrefix: "",
				Site:        site,
				Language:    cfg.DefaultLanguage,
			}

			data := jet.VarMap{}
			data.Set("site", site)

			renderContext(cfg, routeRegistry, ctx, "./", data, w, r)
		}
		return http.HandlerFunc(fn)
	}
}
