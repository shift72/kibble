# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY=kibble

GOARCH = amd64
VERSION=$(shell git describe --abbrev=0 --tags)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -s -w -X github.com/indiereign/shift72-kibble/kibble/version.Version=${VERSION}

.PHONY: clean test

all: test darwin

darwin:
	GOOS=darwin GOARCH=${GOARCH} $(GOINSTALL) -ldflags '${LDFLAGS}' 

test: 
	$(GOTEST) -cover ./...

clean: 
	$(GOCLEAN)

