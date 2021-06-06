
NAME=LeaveApplicationBackend
VERSION=0.1.0
BUILD=`date +%FT%T%z`
# VERSION := $(shell git describe --tags)
# BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

MAKEFLAGS += --silent

BINARY_NOVER=microservice
BINARY=$(PROJECTNAME)_v$(VERSION)

.PHONY: it
it: uninstall clean build install

.PHONY: ci
ci: build

.PHONY: deps
deps:
	go get ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: fmt
fmt: 
	go fmt ./... 	

.PHONY: build
build: fmt deps 
	go build $(LDFLAGS) -o build/$(BINARY)

.PHONY: build-nover
build-nover: fmt deps 
	go build $(LDFLAGS) -o build/$(BINARY_NOVER)	

.PHONY: docker-run
docker-run: clean-nover build-nover 
	docker-compose up -d

.PHONY: clean
clean:
	if [ -f bin/$(BINARY) ] ; then rm build/$(BINARY) ; fi

.PHONY: clean-nover
clean-nover:
	if [ -f bin/$(BINARY_NOVER) ] ; then rm build/$(BINARY_NOVER) ; fi	
