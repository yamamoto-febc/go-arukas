VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
CURRENT_VERSION = $(gobump show -r )
GO_FILES?=$(shell find . -name '*.go')

export GO111MODULE=on

.PHONY: default
default: test

.PHONY: tools
tools:
	go get -u github.com/motemen/gobump

.PHONY: testacc
testacc: 
	TEST_ACC=1 go test ./... $(TESTARGS) -v -timeout=30m -parallel=4 ;

.PHONY: test
test: 
	TEST_ACC=  go test ./... $(TESTARGS) -v -timeout=30m -parallel=4 ;

.PHONY: fmt
fmt:
	gofmt -s -l -w $(GOFMT_FILES)
