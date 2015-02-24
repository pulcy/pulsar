PROJECT := devtool
SCRIPTDIR := $(shell pwd)
VERSION:= $(shell cat VERSION)
COMMIT := $(shell git rev-parse HEAD | cut -c1-8)

GOBUILDDIR := $(SCRIPTDIR)/.gobuild
SRCDIR := $(SCRIPTDIR)/src
BINDIR := $(SCRIPTDIR)

ORGPATH := arvika.subliminl.com/developers
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
	git clone git@github.com:Subliminl/errgo.git $(GOBUILDDIR)/src/github.com/juju/errgo
	git clone git@github.com:Subliminl/go-logging.git $(GOBUILDDIR)/src/github.com/op/go-logging
	git clone git@github.com:Subliminl/pflag.git $(GOBUILDDIR)/src/github.com/spf13/pflag
	git clone git@github.com:Subliminl/cobra.git $(GOBUILDDIR)/src/github.com/spf13/cobra
	git clone git@github.com:Subliminl/go-gitlab-client.git $(GOBUILDDIR)/src/github.com/subliminl/go-gitlab-client
	git clone git@github.com:Subliminl/go-homedir.git $(GOBUILDDIR)/src/github.com/mitchellh/go-homedir
	GOPATH=$(GOPATH) go get github.com/coreos/go-semver/semver

$(BIN): .gobuild $(SOURCES) 
	docker run \
	    --rm \
	    -v $(SCRIPTDIR):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -w /usr/code/src \
	    golang:1.3.1-cross \
	    go build -a -ldflags "-X main.projectVersion $(VERSION) -X main.projectBuild $(COMMIT)" -o ../$(PROJECT)

test:
	#GOPATH=$(GOPATH) go test -v $(REPOPATH)/scheduler
	