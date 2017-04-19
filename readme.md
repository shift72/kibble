# kibble
*def: To chop or grind coarsely*

## Usage
```kibble config``` - configure kibble to use the api key when requesting req

```kibble server --watch``` - starts a server

```kibble render``` - sample render, this is what will be deployed

```kibble help``` - help is here

## Installation
* Install dep ```go get -u github.com/golang/dep/...```
* Install dependencies ```dep ensure -update```
* Build and install ```go install```
* Check installed and running correctly ```kibble version```


## Supports
* model
   * film - done
     * bonus - done
   * taxonomies - done
     * cast - done
     * genre - done
     * year - done
   * pages - done
   * page features - done
   * bundles - done
   * custom pages -> page.html.jet -> page.html - done
     * robots.txt - done
     * humans.txt - done
   * config / toggles - done
   * navigation - done
   * pagination
      * pages - done
* watch files / live reload - done
* cache - done
* support markdown - done
* internationalization - done
   * default language - done
   * language routes /:lang/film/:id - done
* helpers
   * route renders - done
   * canonical route - done
* admin
  * render as admin
  * check admin token is valid
  * request user to (re)login


## TODO
 * model
    * tv -
    * subtitles -
    * pagination
      * language routes
    * features - ?? collections ??
    * json? - what would choice tv tvguide do?
 * shortcodes - these need to happen as part of the markdown process
 * init
    * create a base implementation
 * upload
 * diff / merge ??
