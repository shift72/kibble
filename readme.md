# kibble
*def: To chop or grind coarsely*

## Usage
```kibble config``` - configure kibble to use the api key when requesting req

```kibble render --watch``` - sample render, this is what will be deployed

```kibble help``` - help is here

## Installation
* Requires go 1.8.1 (cgo fix needed)

* Install sass first ```go get github.com/wellington/go-libsass```

* Install dep ```go get -u github.com/golang/dep/...```
* Install dependencies ```dep ensure -update```
* Build and install ```go install```
* Check installed and running correctly ```kibble version```

# LibSass testing
On OSX requires go 1.8.1 which fixes a bug with cgo and xcode
Only appears to build successfully in the default go path
```cd $GOPATH/src/github.com/wellington/go-libsass``
```go test```

# Release management
```go get github.com/mitchellh/gox```
```go get https://github.com/tcnksm/ghr```
gox -osarch="darwin/amd64" -cgo -verbose -rebuild -output="pkg/{{.Dir}}/{.OS}}_{{.Arch}}/kibble"
gox -osarch="linux/amd64" -cgo -verbose -rebuild -output="pkg/{{.Dir}}/{.OS}}_{{.Arch}}/kibble"

# cross compilation ???
brew tap cosmo0920/mingw_w64
brew mingw-w64

/usr/local/bin/i686-w64-mingw32-gcc
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/i686-w64-mingw32-gcc CXX=/usr/local/bin/i686-w64-mingw32-g++ CGO_LDFLAGS="-static" go build -i -v -x -ldflags "-extldflags '-static' -extld=$CC"

GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-ming32-gcc CXX=i686-w64-mingw32-g++ CGO_LDFLAGS="-static" go build -i -v -x -ldflags "-extldflags '-static' -extld=$CC"




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
