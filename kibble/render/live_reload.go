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

package render

import (
	"bytes"
	"fmt"

	"kibble/models"
	"kibble/utils"

	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	fsnotify "gopkg.in/fsnotify/fsnotify.v1"
)

var liveReloadEmbed = `
<script>
(function(){
	var etag = ''; // IE 11 can't handle the 'If-Modified-Since' header
  function checkForChanges() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
			if (this.readyState == this.HEADERS_RECEIVED) {
				if (!etag){
					etag = this.getResponseHeader("Etag");
					return;
				}

				var tag = this.getResponseHeader("Etag");
				if (tag !== etag){
					location.reload(true);
				}
				etag = tag;
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

var kibbleErrorMessage = `
  <div class="kibble-errors">
    <b>Kibble errors</b>
    <p>This may prevent the site from working properly. Check terminal output for details<p>
  </div>

  <script>
    const el = document.querySelector('.kibble-errors');
    const hide = () => {
        el.classList.add("kibble-errors--hidden");
        setTimeout(() => el.remove(), 1000);
    };
    
    el.addEventListener('click', hide);
    setTimeout(hide, 5000);
  </script>
  
  <style>
    .kibble-errors {
        position: fixed;
        bottom: 1rem;
        right: 1rem;
        width: 350px;
        z-index: 10000;
        border-radius: 0.25rem;
        border: 1px solid #a61e4d;
        color: #a61e4d;
        background-color: #fcc2d7;
        box-shadow: 0 10px 10px rgba(0,0,0,0.2);
        transition: transform, opacity, 0.5s;
        transform: translate(0,0);
        opacity: 1;
        font-size: 16px;
        padding: 0.75em 1em;
    }

    .kibble-errors p {
        margin: 0;
        font-size: 0.8em;
        opacity: 0.8;
    }

    .kibble-errors--hidden {
        transform: translateY(100px);
        opacity: 0;
    }
  </style>
`

// LiveReload -
type LiveReload struct {
	lastModified              time.Time
	logReader                 utils.LogReader
	sourcePath                string
	config                    models.LiveReloadConfig
	reloadBrowserOnFileChange bool
}

// WrapperResponseWriter - wraps request
// intercepts all write calls so as to append the live reload script
type WrapperResponseWriter struct {
	http.ResponseWriter
	status       int
	buf          bytes.Buffer
	htmlToInject bytes.Buffer
}

// NewWrapperResponseWriter - create a new response writer
func NewWrapperResponseWriter(w http.ResponseWriter) *WrapperResponseWriter {
	return &WrapperResponseWriter{ResponseWriter: w, status: 200}
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

func (w *WrapperResponseWriter) AddHtmlInject(html string) {
	w.htmlToInject.WriteString(html)
}

// Done - called when are ready to return a result
func (w *WrapperResponseWriter) Done() (err error) {
	if w.htmlToInject.Len() > 0 {
		responseBytes := w.buf.Bytes()

		// Find the end of the HTML body to inject our stuff. This is a little
		// bit flaky as it assumes "sensible" html style on the body closing tag
		endOfBodyIndex := bytes.LastIndex(responseBytes, []byte("</body>"))
		if endOfBodyIndex >= 0 {
			chunks := [][]byte{
				responseBytes[0:endOfBodyIndex],
				w.htmlToInject.Bytes(),
				responseBytes[endOfBodyIndex:],
			}

			length := 0
			for _, chunk := range chunks {
				length += len(chunk)
			}

			w.Header().Set("Content-Length", strconv.Itoa(length))
			w.ResponseWriter.WriteHeader(w.status)

			for _, chunk := range chunks {
				if _, err = w.ResponseWriter.Write(chunk); err != nil {
					return err
				}
			}

			return nil
		}

		// If no end of body can be found, don't inject anything. It might be a
		// partial, or not a valid HTML doc.
	}

	w.Header().Set("Content-Length", strconv.Itoa(w.buf.Len()))
	w.ResponseWriter.WriteHeader(w.status)
	_, err = w.ResponseWriter.Write(w.buf.Bytes())
	return err
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
			strings.HasSuffix(r.RequestURI, ".html") {

			if ww.Status() == 200 {
				if len(live.logReader.Logs()) > 0 {
					ww.AddHtmlInject(kibbleErrorMessage)
				}

				if live.reloadBrowserOnFileChange {
					ww.AddHtmlInject(liveReloadEmbed)
				}
			}
		}

		if err := ww.Done(); err != nil {
			panic(err)
		}
	})
}

// Handler - handle the live reload
func (live *LiveReload) Handler(w http.ResponseWriter, req *http.Request) {
	if !live.reloadBrowserOnFileChange {
		return
	}

	matchEtag := req.Header.Get("If-Modified-Since")

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Etag", live.lastModified.String())

	if matchEtag == live.lastModified.String() || matchEtag == "" {
		w.WriteHeader(http.StatusNotModified)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_, _ = w.Write([]byte(fmt.Sprintf("last modified: %s", live.lastModified.String())))
}

// StartLiveReload - start the process to watch the files and wait for a reload
func (live *LiveReload) StartLiveReload(port int32, fn func()) {

	url := fmt.Sprintf("http://localhost:%d/", port)

	rendered := make(chan bool)

	// wait for changes
	changesChannel := make(chan bool)
	go func() {
		if live.reloadBrowserOnFileChange {
			log.Info("Starting live reload - %s", url)
		} else {
			log.Info("Starting web server - %s", url)
		}

		for range changesChannel {
			now := time.Now()
			// throttle the amount of changes, due to some editors
			// *cough* Sublime Text *cough* sending multiple WRITES for 1 file
			if !live.lastModified.IsZero() && now.Sub(live.lastModified).Seconds() <= 1 {
				log.Debug("Ignoring multiple changes")
				continue
			}

			fn()
			live.lastModified = time.Now()

			// non blocking send
			select {
			case rendered <- true:
			default:
			}
		}
	}()

	if live.config.LaunchBrowser {
		// launch the browser
		go func() {

			// wait for the channel to be rendered
			<-rendered

			utils.LaunchBrowser(url)
		}()
	}

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
					log.Debugf("change (%s) detected: %s", event.Op, event.Name)
					changesChannel <- true
				}
			case err = <-watcher.Errors:
				log.Error("Watcher: ", err)
			}
		}
	}()

	ignorer := utils.NewFileIgnorer(live.sourcePath, live.config.IgnoredPaths)

	// search the path for files that might have changed
	err = filepath.Walk(live.sourcePath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Error("Watcher: ", err)
		}

		if ignorer.IsIgnored(path) {
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
