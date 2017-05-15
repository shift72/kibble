#!/bin/bash
set -e

go test -cover $(go list ./... | grep -v /vendor/)
