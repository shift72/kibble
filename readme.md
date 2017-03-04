==kibble==
To chop or grind coarsely

==dependencies==
go get -v github.com/spf13/cobra/cobra
go get -u golang.org/x/sys/...
go get -u github.com/fsnotify/fsnotify
go get -u github.com/CloudyKit/jet
go get -u github.com/russross/blackfriday
go get -u github.com/microcosm-cc/bluemonday
go get -u github.com/nicksnyder/go-i18n/goi18n
go get -u github.com/peterbourgon/diskv
go get -u github.com/gregjones/httpcache
go get -u github.com/gosimple/slug

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
    * config / toggles - done
    * navigation - done
 * populate model from api
 * shortcodes - these need to happen as part of the markdown process
 * live reload
 * init
    * create a base implementation
 * upload
 * diff / merge ??
 * watch files / live reload
 * admin - unpublished / published
 * cache - done
 * support markdown - done
 * internationalization - done
    * default language - done
    * language routes /:lang/film/:id - done
 * helpers
    * route renders - done
