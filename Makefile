BUILD_BIN_PATH := $(shell pwd)/bin
BUILD_BIN_NAME := beat-exporter

VERSION_PKG := beat-exporter/internal/service

LD_FLAGS += -X "$(VERSION_PKG).releaseVersion=$(shell git describe --tags --dirty --always)"
LD_FLAGS += -X "$(VERSION_PKG).buildDate=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LD_FLAGS += -X "$(VERSION_PKG).gitHash=$(shell git rev-parse HEAD)"
LD_FLAGS += -X "$(VERSION_PKG).gitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

GOVER_MAJOR := $(shell go version | sed -E -e "s/.*go([0-9]+)[.]([0-9]+).*/\1/")
GOVER_MINOR := $(shell go version | sed -E -e "s/.*go([0-9]+)[.]([0-9]+).*/\2/")
GO111 := $(shell [ $(GOVER_MAJOR) -gt 1 ] || [ $(GOVER_MAJOR) -eq 1 ] && [ $(GOVER_MINOR) -ge 11 ]; echo $$?)
ifeq ($(GO111), 1)
  $(error "go below 1.11 does not support modules")
endif

default: check build

build-deps:
	@mkdir -p bin
build: export GO111MODULE=on
build: build-deps
	go build -ldflags '$(LD_FLAGS)' -o $(BUILD_BIN_PATH)/$(BUILD_BIN_NAME) main.go

check: check-static
check-static: tools/bin/golangci-lint
	tools/bin/golangci-lint run -v --deadline=3m --config ./.golangci.yaml $$($(PACKAGE_DIRECTORIES))
tools/bin/golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./tools/bin v1.41.1;