==kibble==
To chop or grind coarsely

==dependencies==
go get -v github.com/spf13/cobra/cobra
go get -u golang.org/x/sys/...
go get github.com/fsnotify/fsnotify
go get -u github.com/CloudyKit/jet

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
    * json? - what would choice tv tvguide do?
    * config / toggles
 * populate model from api
 * cache
 * live reload
 * internationalization
    * default language
    * language routes /:lang/film/:id
 * helpers
    * route renders
 * upload
 * diff / merge ??
 * watch files / live reload
