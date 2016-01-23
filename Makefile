PROJECT := pulcy
SCRIPTDIR := $(shell pwd)
VERSION:= $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

GOBUILDDIR := $(SCRIPTDIR)/.gobuild
SRCDIR := $(SCRIPTDIR)/src
BINDIR := $(SCRIPTDIR)

ORGPATH := git.pulcy.com/pulcy
ORGDIR := $(GOBUILDDIR)/src/$(ORGPATH)
REPONAME := $(PROJECT)
REPODIR := $(ORGDIR)/$(REPONAME)
REPOPATH := $(ORGPATH)/$(REPONAME)
BIN := $(BINDIR)/$(PROJECT)

GOPATH := $(GOBUILDDIR)

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

.gobuild:
	mkdir -p $(ORGDIR)
	rm -f $(REPODIR) && ln -s ../../../../src $(REPODIR)
	git clone git@github.com:juju/errgo.git $(GOBUILDDIR)/src/github.com/juju/errgo
	git clone git@github.com:op/go-logging.git $(GOBUILDDIR)/src/github.com/op/go-logging
	git clone git@github.com:spf13/pflag.git $(GOBUILDDIR)/src/github.com/spf13/pflag
	git clone git@github.com:spf13/cobra.git $(GOBUILDDIR)/src/github.com/spf13/cobra
	git clone git@github.com:inconshreveable/mousetrap.git $(GOBUILDDIR)/src/github.com/inconshreveable/mousetrap
	git clone git@github.com:cpuguy83/go-md2man.git $(GOBUILDDIR)/src/github.com/cpuguy83/go-md2man
	git clone git@github.com:russross/blackfriday.git $(GOBUILDDIR)/src/github.com/russross/blackfriday
	git clone git@github.com:shurcooL/sanitized_anchor_name.git $(GOBUILDDIR)/src/github.com/shurcooL/sanitized_anchor_name
	git clone git@github.com:ewoutp/go-gitlab-client.git $(GOBUILDDIR)/src/github.com/ewoutp/go-gitlab-client
	git clone git@github.com:mitchellh/go-homedir.git $(GOBUILDDIR)/src/github.com/mitchellh/go-homedir
	GOPATH=$(GOPATH) go get github.com/coreos/go-semver/semver
	GOPATH=$(GOPATH) go get github.com/mgutz/ansi

$(BIN): .gobuild $(SOURCES)
	cd $(SRCDIR) &&    go build -a -ldflags "-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" -o ../$(PROJECT)

test:
	#GOPATH=$(GOPATH) go test -v $(REPOPATH)/scheduler
