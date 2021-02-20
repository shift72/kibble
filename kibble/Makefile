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
LDFLAGS = -s -w -X kibble/version.Version=${VERSION}

.PHONY: clean test

all: test darwin

darwin:
	GOOS=darwin GOARCH=${GOARCH} $(GOBUILD) -ldflags '${LDFLAGS}'

linux:
	GOOS=linux GOARCH=${GOARCH} $(GOBUILD) -ldflags '${LDFLAGS}' -o ${BINARY} .

linux-install:
	GOOS=linux GOARCH=${GOARCH} $(GOINSTALL) -ldflags '${LDFLAGS}' -v .

windows:
	GOOS=windows GOARCH=${GOARCH} $(GOBUILD) -ldflags '${LDFLAGS}' -o ${BINARY}-windows-${GOARCH}.exe .

test:
	$(GOTEST) -cover ./...

run-darwin: darwin
	cd sample_site && ../kibble render

run-linux: linux
	cd sample_site && ../kibble render -v

clean:
	$(GOCLEAN)

