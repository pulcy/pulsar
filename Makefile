PROJECT := subliminl
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

.PHONY: clean test

all: $(BIN)

clean:
	rm -Rf $(BIN) $(GOBUILDDIR)

.gobuild: 
	mkdir -p $(ORGDIR)
	ln -s $(SRCDIR) $(REPODIR)
	git clone git@github.com:Subliminl/errgo.git $(GOBUILDDIR)/src/github.com/juju/errgo
	git clone git@github.com:Subliminl/go-logging.git $(GOBUILDDIR)/src/github.com/op/go-logging
	git clone git@github.com:Subliminl/pflag.git $(GOBUILDDIR)/src/github.com/spf13/pflag
	git clone git@github.com:Subliminl/cobra.git $(GOBUILDDIR)/src/github.com/spf13/cobra
	git clone git@github.com:Subliminl/semver.git $(GOBUILDDIR)/src/github.com/blang/semver

$(BIN): .gobuild $(SOURCES) 
	GOBIN=$(BINDIR) GOPATH=$(GOPATH) go build -ldflags "-X main.projectVersion $(VERSION) -X main.projectBuild $(COMMIT)" -o $(BIN) $(REPOPATH)

test:
	#GOPATH=$(GOPATH) go test -v $(REPOPATH)/scheduler
	