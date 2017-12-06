# kibble
*def: To chop or grind coarsely*


```npm install -g kibble``` - installs kibble

```kibble init``` - find a template to start with

```cd new-template``` 

```npm start``` - starts kibble

```npm run publish``` - publish to the current site


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

# Publish
Publish will zip all files placed in the ```/.kibble/dist``` directory

```
/.kibble
  /dist       - publish directory
  /kibble.zip - zip file to be published
```

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
   * subtitles - done
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
* shortcodes 
  * common shortcodes - done
* sync
  * upload - done

## TODO
 * model
    * tv - 
    * json? - what would choice tv tvguide do?
 
 * download
    * build and release on github   
    
 * upload
  * zip - https://godoc.org/github.com/pierrre/archivefile/zip
  * build process ?? sass / less
  * api to upload the file

 * diff / merge ??
