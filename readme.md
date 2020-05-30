# kibble

def: To chop or grind coarsely

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

* Requires go 1.13.0 (go mod) 
* Build and install ```go install```
* Check installed and running correctly ```kibble version```

## Publish

Publish will zip all files placed in the ```/.kibble/dist``` directory

```bash

/.kibble
  /dist       - publish directory
  /kibble.zip - zip file to be published

```

## Releasing new versions

Kibble is released to 3 places

  1. github - mostly historical purposes, intention is to make this open source
  2. SHIFT72 Platform - this is where the platform will pull the kibble release from
  3. NPM - to support installation for third parties via npm and the node js environment

``` bash

# a. update npm package version
#    kibble-npm/package.json

# b. ensure that the release is tagged correctly. Miss the prepended 'v' as this will mess S3 up.
git tag 0.9.6 master
git push origin 0.9.6

# c. build and release to locations 1 and 2
cd kibble
goreleaser --rm-dist

# d. release to location 3
cd kibble-npm
npm publish

# e. register new build in uber admin

# f. update any sample templates with the new kibble version

```

## Supports

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
