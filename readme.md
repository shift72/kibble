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


Supports
* model
   * film - done
     * bonus - done
   * pages - done
   * page features - done
   * bundles - done
   * custom pages -> page.html.jet -> page.html - done
     * robots.txt - done
     * humans.txt - done
   * config / toggles - done
   * navigation - done
* watch files / live reload - done
* cache - done
* support markdown - done
* internationalization - done
   * default language - done
   * language routes /:lang/film/:id - done
* helpers
   * route renders - done
   * canonical route - done


TODO:
 * model
    * tv -
    * subtitles -
    * pagination - ??
    * features - ?? collections ??
    * taxonomies -
      * cast -
      * genre -
      * year -
    * json? - what would choice tv tvguide do?
 * shortcodes - these need to happen as part of the markdown process
 * init
    * create a base implementation
 * upload
 * diff / merge ??
 * admin
    * unpublished / published
    * requesting using authtoken
