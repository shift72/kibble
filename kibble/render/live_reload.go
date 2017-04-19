package render

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	fsnotify "gopkg.in/fsnotify.v1"
)

var embed = `
<script>
(function(){
  function checkForChanges() {
    xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4 && this.status == 200) {
        location.reload(true);
      }
    };
    xhttp.open("GET", "/kibble/live_reload", true);
    xhttp.send();
    setTimeout(checkForChanges, 3000);
  }
  setTimeout(checkForChanges, 3000);
})();
</script>
`

// LiveReload -
type LiveReload struct {
	changed bool
}

var ignorePaths = []string{".git", ".kibble"}

// GetMiddleware - return a handler
func (live *LiveReload) GetMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			w.Write([]byte(embed))
		}
		return http.HandlerFunc(fn)
	}
}

// Handler - handle the live reload
func (live *LiveReload) Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	if live.changed {
		live.changed = false
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotModified)
	}
	w.Write([]byte("kibble...\r\n"))
}

// StartLiveReload - start the process to watch the files and wait for a reload
func (live *LiveReload) StartLiveReload(fn func()) {

	// wait for changes
	changesChannel := make(chan bool)
	go func() {
		fmt.Println("starting live reload")
		for _ = range changesChannel {
			fn()
			live.changed = true
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
