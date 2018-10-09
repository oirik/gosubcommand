NAME      := $(shell basename `pwd`)
VERSION   := $(shell git describe --tags --abbrev=0)
REVISION  := $(shell git rev-parse --short HEAD)
GOLINT    := $(shell command -v golint 2> /dev/null)
LDFLAGS   := -X 'main.Version=$(VERSION)' -X 'main.Revision=$(REVISION)'

.PHONY: test
test: lint
	go test -race -v ./...

.PHONY: golint
golint:
ifndef GOLINT
	go get -u github.com/golang/lint/golint
endif

.PHONY: clean
clean:
	go clean

.PHONY: lint
lint: golint deps
	go vet ./...
	golint -set_exit_status ./...
