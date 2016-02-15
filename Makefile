PROJECT := pulcy
SCRIPTDIR := $(shell pwd)
VERSION:= $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

GOBUILDDIR := $(SCRIPTDIR)/.gobuild
SRCDIR := $(SCRIPTDIR)
BINDIR := $(SCRIPTDIR)
VENDORDIR = $(SCRIPTDIR)/vendor

ORGPATH := github.com/pulcy
ORGDIR := $(GOBUILDDIR)/src/$(ORGPATH)
REPONAME := $(PROJECT)
REPODIR := $(ORGDIR)/$(REPONAME)
REPOPATH := $(ORGPATH)/$(REPONAME)
BIN := $(BINDIR)/$(PROJECT)

GOPATH := $(GOBUILDDIR)
GOVERSION := 1.5.3

SOURCES := $(shell find $(SRCDIR) -name '*.go')

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif
ifndef GOARCH
	GOARCH := $(shell go env GOARCH)
endif


.PHONY: clean test

all: $(BIN)

clean:
	rm -Rf $(BIN) $(GOBUILDDIR)

$(GOBUILDDIR):
	mkdir -p $(ORGDIR)
	rm -f $(REPODIR) && ln -s ../../../../ $(REPODIR)

update-vendor:
	rm -Rf $(VENDORDIR)
	$(BIN) go vendor -V $(VENDORDIR) \
		github.com/coreos/go-semver/semver \
		github.com/cpuguy83/go-md2man \
		github.com/ewoutp/go-gitlab-client \
		github.com/inconshreveable/mousetrap \
		github.com/juju/errgo \
		github.com/mgutz/ansi \
		github.com/mitchellh/go-homedir \
		github.com/op/go-logging \
		github.com/russross/blackfriday \
		github.com/shurcooL/sanitized_anchor_name \
		github.com/spf13/cobra \
		github.com/spf13/pflag

$(BIN): $(GOBUILDDIR) $(SOURCES)
	docker run \
	    --rm \
	    -v $(SRCDIR):/usr/code \
	    -e GO15VENDOREXPERIMENT=1 \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -w /usr/code/ \
	    golang:$(GOVERSION) \
	    go build -a -ldflags "-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" -o /usr/code/$(PROJECT) $(REPOPATH)

test:
	#GOPATH=$(GOPATH) go test -v $(REPOPATH)/scheduler
