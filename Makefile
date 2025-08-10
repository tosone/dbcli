NAME             := dbcli
ROOT             := github.com/tosone/$(NAME)
GOOS             ?= linux
GOARCH           ?= amd64
PLATFORM         ?= $(GOOS)/$(GOARCH)
CMD_DIR          := .
OUTPUT_DIR       := ./bin
BUILD_DIR        := ./build
DOCKER           ?= docker
VERSION          ?= $(shell git describe --tags --always --dirty)
GITCOMMIT        ?= $(shell git rev-parse HEAD)
GITTREESTATE     ?= $(if $(shell git status --porcelain),dirty,clean)
BUILDDATE        ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

GOLDFLAGS        += -X $(ROOT)/pkg/version.Version=$(shell git describe --tags --always)
GOLDFLAGS        += -X $(ROOT)/pkg/version.BuildDate=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOLDFLAGS        += -X $(ROOT)/pkg/version.GitHash=$(shell git rev-parse --short HEAD)
GOFLAGS           = -ldflags '-s -w $(GOLDFLAGS)' -trimpath

.PHONY: build
build:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o $(OUTPUT_DIR)/$(NAME) $(GOFLAGS) $(CMD_DIR)
