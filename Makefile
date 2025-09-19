GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

.PHONY: all build clean test

all: build

build:
	$(GOBUILD) -o bin/inline-example ./examples/inline

clean:
	$(GOCLEAN)
	rm -f bin/inline-example

test:
	$(GOTEST) ./...


