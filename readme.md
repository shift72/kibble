# kibble

> def: To chop or grind coarsely

[Changelog](changelog.md)

```npm install -g kibble``` - installs kibble

```kibble init``` - find a template to start with

```cd new-template```

```npm start``` - starts kibble

```npm run publish``` - publish to the current site


## Documentation

Documentation is available [on the wiki](https://github.com/shift72/kibble/wiki)

## Usage

```kibble config``` - configure kibble to use the api key when requesting req

```kibble render --watch``` - sample render, this is what will be deployed

```kibble sync``` - sync files to a remote s3 bucket

```kibble help``` - help is here

## Installation

* Requires go 1.17.0 (go mod)
* Build and install ```go install```
* Check installed and running correctly ```kibble version```


## Release process

See [release.md](./release.md) for details.