==kibble==
To chop or grind coarsely

==dependencies==
go get -v github.com/spf13/cobra/cobra
go get -u golang.org/x/sys/...
go get github.com/fsnotify/fsnotify
go get -u github.com/CloudyKit/jet
go get -u github.com/russross/blackfriday
go get -u github.com/microcosm-cc/bluemonday
go get -u github.com/nicksnyder/go-i18n/goi18n

Commands
 - config - set the API_KEY??
 - render - renders the entire site
 - serve - starts a web server for local development

TODO:
 * model
    * film
    * tv
    * pages
    * collections
    * taxonomies
      * cast
      * genre
      * year
    * json? - what would choice tv tvguide do?
    * config / toggles
 * populate model from api
 * cache
 * support markdown - done
 * shortcodes - these need to happen as part of the markdown process
 * live reload
 * admin - unpublished
 * internationalization - done
    * default language - done
    * language routes /:lang/film/:id - done
 * helpers
    * route renders - done
 * upload
 * diff / merge ??
 * watch files / live reload
