#!/bin/bash
set -e

go test $(go list ./... | grep -v /vendor/)
