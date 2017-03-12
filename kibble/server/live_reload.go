package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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

// useful to trigger one new reload if the site has turned live reload off
var changed = true

func startLiveReload() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					changed = true
				}
			case err = <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	var searchDir = "./templates"

	err = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
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

// InjectLiveReloadScript - middleware to append script to check for changes
func InjectLiveReloadScript(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		w.Write([]byte(embed))
	}
	return http.HandlerFunc(fn)
}

func handleLiveReload(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	if changed {
		changed = false
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotModified)
	}
	w.Write([]byte("kibble...\r\n"))
}
