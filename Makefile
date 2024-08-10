PROJECT := pulsar
SCRIPTDIR := $(shell pwd)
VERSION:= $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

GOBUILDDIR := $(SCRIPTDIR)/.gobuild
SRCDIR := $(SCRIPTDIR)
BINDIR := $(SCRIPTDIR)/bin
VENDORDIR = $(SCRIPTDIR)/vendor

ORGPATH := github.com/pulcy
ORGDIR := $(GOBUILDDIR)/src/$(ORGPATH)
REPONAME := $(PROJECT)
REPODIR := $(ORGDIR)/$(REPONAME)
REPOPATH := $(ORGPATH)/$(REPONAME)
BIN := $(BINDIR)/$(PROJECT)

SOURCES := $(shell find $(SRCDIR) -name '*.go')

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif
ifndef GOARCH
	GOARCH := $(shell go env GOARCH)
endif


.PHONY: clean binaries test release tgz

all: binaries

clean:
	rm -Rf $(BIN) bin

binaries: $(SOURCES)
	mkdir -p bin
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64" \
		-ldflags="-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" \
		-output="bin/{{.OS}}/{{.Arch}}/$(PROJECT)" \
		-tags="netgo" \
		./...

release: tgz

tgz: binaries
	tar zcf bin/$(PROJECT)-linux-amd64.tgz -C $(BINDIR)/linux/amd64 $(PROJECT)
	tar zcf bin/$(PROJECT)-linux-arm64.tgz -C $(BINDIR)/linux/arm64 $(PROJECT)
	tar zcf bin/$(PROJECT)-darwin-amd64.tgz -C $(BINDIR)/darwin/amd64 $(PROJECT)
	tar zcf bin/$(PROJECT)-darwin-arm64.tgz -C $(BINDIR)/darwin/arm64 $(PROJECT)

test:
	#GOPATH=$(GOPATH) go test -v $(REPOPATH)/scheduler
