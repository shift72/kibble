# kibble
*def: To chop or grind coarsely*

## Usage
```kibble config``` - configure kibble to use the api key when requesting req

```kibble render --watch``` - sample render, this is what will be deployed

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
  * render as admin - done
  * check admin token is valid - done
  * request user to (re)login - done
* init
  * create a base implementation - done
  * find repo based on - done
  * clone the repo - done

## TODO
 * model
    * tv -
    * subtitles -
    * pagination
      * language routes
    * features - ?? collections ??
    * json? - what would choice tv tvguide do?
 * shortcodes - these need to happen as part of the markdown process
 * download
    * build and release on github
    *
    *
    
 * upload
  * zip - https://godoc.org/github.com/pierrre/archivefile/zip
  * build process ?? sass / less
  * api to upload the file

 * diff / merge ??
