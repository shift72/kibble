# kibble
*def: To chop or grind coarsely*

## Usage
```kibble config``` - configure kibble to use the api key when requesting req

```kibble render --watch``` - sample render, this is what will be deployed

```kibble sync``` - sync files to a remote s3 bucket

```kibble help``` - help is here

## Installation
* Requires go 1.8.1 (cgo fix needed)
* Install dep ```go get -u github.com/golang/dep/...```
* Install dependencies ```dep ensure -update```. This might not work, so try ```dep ensure -v github.com/fsnotify/fsnotify@^1.4.2``` in case of errors.
* Build and install ```go install```
* Check installed and running correctly ```kibble version```

# LibSass Build
TODO:

# build and build order
/static           >  copy  > .kibble/
/styles           >  sass  > .kibble/styles/main.css
/templates        > render > .kibble/
[**/*.jet]        > render > .kibble/
[languages.json]  > render > .kibble/
[site.json]       

# zip structure
/static           >        zip > /static
/styles           > sass > zip > /styles/main.css
/templates        >        zip > /templates
[languages.json]  >        zip > /
[**/*.jet]        >        zip > /
[site.json]       >        zip > /

# server side render
/static           >  copy  > .kibble/
/styles           >  sass  > .kibble/styles/main.css
/templates        > render > .kibble/
[**/*.jet]        > render > .kibble/
[languages.json]  > render > .kibble/
[site.json]       


## Supports
* model
   * film - done
     * bonus - done
   * taxonomies - done
     * cast - done
     * genre - done
     * year - done
   * pages - done
     * type - templates - done
   * page features - done
   * bundles - done
   * custom pages -> page.html.jet -> page.html - done
     * robots.txt - done
     * humans.txt - done
   * config / toggles - done
   * navigation - done
   * pagination
      * pages - done
      * language routes - done
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
* sync
  * 

## TODO
 * model
    * tv -
    * subtitles -
    
    * json? - what would choice tv tvguide do?
 * shortcodes - these need to happen as part of the markdown process
 * download
    * build and release on github   
    
 * upload
  * zip - https://godoc.org/github.com/pierrre/archivefile/zip
  * build process ?? sass / less
  * api to upload the file

 * diff / merge ??
