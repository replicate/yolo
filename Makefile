SHELL := /bin/bash

DESTDIR ?=
PREFIX = /usr/local
BINDIR = $(PREFIX)/bin

INSTALL := install -m 0755
INSTALL_PROGRAM := $(INSTALL)

GO := go
GOOS := $(shell $(GO) env GOOS)
GOARCH := $(shell $(GO) env GOARCH)

default: all

.PHONY: all
all: yolo


.PHONY: yolo
yolo:
	$(eval yolo_VERSION ?= $(shell git describe --tags --match 'v*' --abbrev=0)+dev)
	CGO_ENABLED=0 $(GO) build -o $@ \
		-ldflags "-X github.com/replicate/yolo/pkg/global.BuildTime=$(shell date +%Y-%m-%dT%H:%M:%S%z) -w" \
		yolo.go

.PHONY: install
install: yolo
	$(INSTALL_PROGRAM) -d $(DESTDIR)$(BINDIR)
	$(INSTALL_PROGRAM) yolo $(DESTDIR)$(BINDIR)/yolo

.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(BINDIR)/yolo

.PHONY: clean
clean:
	$(GO) clean
	rm -rf build dist
	rm -f yolo

.PHONY: test-go
test-go: vet lint-go

.PHONY: test-integration
test-integration: yolo
	cd test-integration/ && $(MAKE) PATH="$(PWD):$(PATH)" test


.PHONY: test
test: test-go


.PHONY: generate
generate:
	$(GO) generate ./...


.PHONY: vet
vet:
	$(GO) vet ./...


.PHONY: check-fmt
check-fmt:
	$(GO) run golang.org/x/tools/cmd/goimports -d .
	@test -z $$($(GO) run golang.org/x/tools/cmd/goimports -l .)

.PHONY: lint-go
lint-go:
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...


.PHONY: lint
lint: lint-go lint-python

.PHONY: mod-tidy
mod-tidy:
	$(GO) mod tidy