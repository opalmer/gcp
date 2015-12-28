GO ?= $(shell which go)
GOLINT ?= $(GOPATH)/bin/golint


build: fmt lint
	@echo "go build"
	@$(GO) build

dep:
	@echo "Retrieving dependencies..."
	@$(GO) get
	@$(GO) get -u github.com/golang/lint/golint

fmt:
	@echo "go fmt"
	@$(GO) fmt main.go
	@find db -type f -name "*.go" -exec $(GO) fmt {} \;
	@find server -type f -name "*.go" -exec $(GO) fmt {} \;

lint:
	@echo "golint"
	@$(GOLINT) main.go
	@$(GOLINT) db
	@$(GOLINT) server

server: build
	@./goback -server