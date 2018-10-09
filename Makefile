NAME      := $(shell basename `pwd`)
REVISION  := $(shell git rev-parse --short HEAD)
GODEP     := $(shell command -v dep 2> /dev/null)
GOLINT    := $(shell command -v golint 2> /dev/null)
LDFLAGS   := -X 'main.Version=$(VERSION)' -X 'main.Revision=$(REVISION)'
VENDORDIR :=./vendor

.PHONY: test
test: lint
	go test -race -v ./...

.PHONY: godep
godep:
ifndef GODEP
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: golint
golint:
ifndef GOLINT
	go get -u github.com/golang/lint/golint
endif

.PHONY: deps
deps: godep
	dep ensure

.PHONY: clean
clean:
	go clean
	rm -rf $(VENDORDIR)/*

.PHONY: lint
lint: golint deps
	go vet ./...
	golint -set_exit_status ./

