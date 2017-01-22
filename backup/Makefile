COMMIT     := $(shell git rev-parse --short HEAD)
VERSION    := 1.0.1

LDFLAGS    := -ldflags \
              "-X main.Commit $(COMMIT)\
               -X main.Version $(VERSION)"

GOOS       := $(shell go env GOOS)
GOARCH     := $(shell go env GOARCH)
GOBUILD    := GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o dist/backups $(LDFLAGS)
GOFILES    := $(shell find . -name "*.go" -exec echo {}  \; | sed -e s/.\\/// | grep -ve test)

ARCHIVE    := backups-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz
DISTDIR    := dist/$(GOOS)_$(GOARCH)

.PHONY: default archive clean install

default: *.go
	$(GOBUILD)

archive: dist/$(ARCHIVE)

GO ?= go

all: build

build:
	$(GO) build

check:
	@$(GO) list -f '{{join .Deps "\n"}}' | xargs $(GO) list -f '{{if not .Standard}}{{.ImportPath}} {{.Dir}}{{end}}' | column -t

clean:
	git clean -f -x -d

run:
	@$(GO) run $(GOFILES) -http.addr=:8080

dist/$(ARCHIVE): $(DISTDIR)/api
	tar -C $(DISTDIR) -czvf $@ .

$(DISTDIR)/api: *.go
	$(GOBUILD) -o $@
