package render

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var probePaths = []string{"", "index.html"}

// StaticMiddleware -
func StaticMiddleware() func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			foundPath := ""

			// check if the request + jet file exists
			for _, probe := range probePaths {
				path := fmt.Sprintf("%s%s%s", rootPath, r.RequestURI, probe)
				stat, err := os.Stat(path)

				if stat != nil && stat.IsDir() {
					continue
				}

				if !os.IsNotExist(err) {
					foundPath = path
					break
				}
			}

			if foundPath == "" {
				next.ServeHTTP(w, r)
				return
			}

			b, err := ioutil.ReadFile(foundPath)
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
