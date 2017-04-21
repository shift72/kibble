package render

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	fsnotify "gopkg.in/fsnotify.v1"
)

var embed = `
<script>
(function(){
	var etag = '';
  function checkForChanges() {
    xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
			if (this.readyState == this.HEADERS_RECEIVED) {
        etag = this.getResponseHeader("Etag");
      }
      if (this.readyState == 4 && this.status == 200) {
        location.reload(true);
      }
    };
    xhttp.open("GET", "/kibble/live_reload", true);
		xhttp.setRequestHeader("If-Modified-Since", etag);
    xhttp.send();
    setTimeout(checkForChanges, 3000);
  }
  setTimeout(checkForChanges, 3000);
})();
</script>
`

// LiveReload -
type LiveReload struct {
	lastModified time.Time
}

// WrapperResponseWriter - used to track the status of a response
type WrapperResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func NewWrapperResponseWriter(w http.ResponseWriter) *WrapperResponseWriter {
	return &WrapperResponseWriter{ResponseWriter: w}
}

func (w *WrapperResponseWriter) Status() int {
	return w.status
}

func (w *WrapperResponseWriter) Write(p []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (w *WrapperResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	// Check after in case there's error handling in the wrapped ResponseWriter.
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
}

var ignorePaths = []string{".git", ".kibble"}

// GetMiddleware - return a handler
func (live *LiveReload) GetMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := NewWrapperResponseWriter(w)
			next.ServeHTTP(ww, r)

			if ww.Status() == 200 {
				ww.Write([]byte(embed))
			}
		}
		return http.HandlerFunc(fn)
	}
}

// Handler - handle the live reload
func (live *LiveReload) Handler(w http.ResponseWriter, req *http.Request) {
	matchEtag := req.Header.Get("If-Modified-Since")

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Etag", live.lastModified.String())

	if matchEtag == live.lastModified.String() || matchEtag == "" {
		w.WriteHeader(http.StatusNotModified)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Write([]byte(fmt.Sprintf("last modified: %s", live.lastModified.String())))
}

// StartLiveReload - start the process to watch the files and wait for a reload
func (live *LiveReload) StartLiveReload(fn func()) {

	// wait for changes
	changesChannel := make(chan bool)
	go func() {
		fmt.Println("starting live reload")
		for _ = range changesChannel {
			fn()
			live.lastModified = time.Now()
		}
	}()

	// useful to trigger one new reload
	changesChannel <- true

	live.selectFilesToWatch(changesChannel)
}

func (live *LiveReload) selectFilesToWatch(changesChannel chan bool) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// listen for fs events and pass via channel
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					changesChannel <- true
				}
			case err = <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// search the path for files that might have changed
	var searchDir = "."
	err = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if ignorePath(path) {
			return nil
		}

		if f.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				log.Fatal("unable to watch dir", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func ignorePath(name string) bool {
	for _, c := range ignorePaths {
		if strings.HasPrefix(name, c) {
			return true
		}
	}
	return false
}
