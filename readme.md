# kibble

> def: To chop or grind coarsely

[Changelog](changelog.md)

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

* Requires go 1.17.0 (go mod)
* Build and install ```go install```
* Check installed and running correctly ```kibble version```

# Supports

* model
  * film - done
    * bonus - done
  * tv - done
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
  * return errors via summary
* download
  * build and release on github
  * deploy via npm