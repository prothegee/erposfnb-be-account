SHELL := /bin/bash

.PHONY: help
help:
	@echo "make help			show this help"
	@echo "make test			run the whole test"
	@echo "make test-verbose	run the whole test with verbose log"
	@echo "make fmt				format whole go files"
	@echo "make format			format whole go files"

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: format
format:
	go fmt ./...
