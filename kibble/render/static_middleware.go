package render

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// StaticMiddleware -
func StaticMiddleware() func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			path := path.Join(rootPath, r.RequestURI)

			// check if the request + jet file exists
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				next.ServeHTTP(w, r)
				return
			}

			b, err := ioutil.ReadFile(path)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(b)
		}
		return http.HandlerFunc(fn)
	}
}
