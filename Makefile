SHELL = /bin/bash
.DEFAULT_GOAL := all

APP := match-maker

BIN_DIR := bin
DIST_DIR := dist

GO_FILES = $(shell find . -type f -name "*.go"  -not -name "mock_*.go" -not -name "*_test.go" $(exclude_dirs))
GO_TEST_COVERAGE_PATH := coverage.out
PACKAGE := github.com/TanyEm/match-maker

all: install generate build

install: go.mod go.sum ## Install dependencies
	go mod download
	go mod tidy -v
.PHONY: install

generate:
	go get go.uber.org/mock/mockgen@v0.5.0
	go generate ./...
	go mod tidy -v
.PHONY: generate

.PHONY: test
test: install ## Run tests.
	go test ./... -coverprofile cover.out

.PHONY: test-verbose
test-verbose: install ## Run tests.
	go test ./... -coverprofile cover.out -v

build: $(BIN_DIR)/$(APP) ## Compile app for local to bin/

$(BIN_DIR)/%: $(GO_FILES) install
	go build -o $@ ./cmd/$*

dist: $(DIST_DIR)/$(APP) ## Compile app for Docker to dist/

$(DIST_DIR)/%: $(GO_FILES)
# Compilation options:
#   - GOOS=linux: explicitly target Linux
#   - GOARCH: explicitly target 64bit CPU
#   - -trimpath: improve reproducibility by trimming the pwd from the binary
#   - -ldflags: extra linker flags
#     - -s: omit the symbol table and debug information making the binary smaller
#     - -w: omit the DWARF symbol table making the binary smaller
#   - -tags: extra tags
#     - osusergo: Use native Go os/user package instead of using OS libs
#     - netgo: Use native Go net package instead of using OS and C libs
#     - static|static_all|static_build: left out because of no clear documentation on them
	GOOS=linux GOARCH=amd64 \
	go build \
		-trimpath \
		-ldflags "-s -w" \
		-tags "osusergo netgo" \
		-o $@ \
		./cmd/$*
