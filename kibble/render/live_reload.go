package render

import (
	"bytes"
	"fmt"

	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

var ignorePaths = []string{".git", ".kibble"}

// LiveReload -
type LiveReload struct {
	lastModified time.Time
}

// WrapperResponseWriter - wraps request
// intercepts all write calls so as to append the live reload script
type WrapperResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	buf         bytes.Buffer
}

// NewWrapperResponseWriter - create a new response writer
func NewWrapperResponseWriter(w http.ResponseWriter) *WrapperResponseWriter {
	return &WrapperResponseWriter{ResponseWriter: w}
}

// Status - get the status
func (w *WrapperResponseWriter) Status() int {
	return w.status
}

// Write - wrap the write
func (w *WrapperResponseWriter) Write(p []byte) (n int, err error) {
	w.buf.Write(p)
	return len(p), nil
}

// Done - called when are ready to return a result
func (w *WrapperResponseWriter) Done() (n int, err error) {
	w.Header().Set("Content-Length", strconv.Itoa(w.buf.Len()))
	w.ResponseWriter.WriteHeader(w.status)
	return w.ResponseWriter.Write(w.buf.Bytes())
}

// WriteHeader - wrap the write header
func (w *WrapperResponseWriter) WriteHeader(code int) {
	w.status = code
}

// GetMiddleware - return a handler
func (live *LiveReload) GetMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ww := NewWrapperResponseWriter(w)
		next.ServeHTTP(ww, r)

		if strings.HasSuffix(r.RequestURI, "/") ||
			strings.HasSuffix(r.RequestURI, "/index.html") {

			if ww.Status() == 200 {
				ww.Write([]byte(embed))
			}
		}

		ww.Done()
	})
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
func (live *LiveReload) StartLiveReload(port int32, fn func()) {

	rendered := make(chan bool)

	// wait for changes
	changesChannel := make(chan bool)
	go func() {
		log.Info("Starting live reload")

		for _ = range changesChannel {
			fn()
			live.lastModified = time.Now()

			// non blocking send
			select {
			case rendered <- true:
			default:
			}
		}
	}()

	// launch the browser
	go func() {

		// wait for the channel to be rendered
		<-rendered

		cmd := exec.Command("open", fmt.Sprintf("http://localhost:%d/", port))
		err := cmd.Start()
		if err != nil {
			log.Error("Watcher: ", err)
		}
	}()

	// useful to trigger one new reload
	changesChannel <- true

	live.selectFilesToWatch(changesChannel)
}

func (live *LiveReload) selectFilesToWatch(changesChannel chan bool) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		//	log.Fatal(err)
	}

	// listen for fs events and pass via channel
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debugf("change detected:", event.Name)
					changesChannel <- true
				}
			case err = <-watcher.Errors:
				log.Error("Watcher: ", err)
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
				log.Error("Watcher: ", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Error("Watcher: ", err)
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
