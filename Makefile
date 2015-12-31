GO ?= go
GOLINT ?= golint

build: fmt lint
	$(GO) build gcp.go

dep:
	$(GO) get
	$(GO) get -u github.com/golang/lint/golint

fmt:
	$(GO) fmt libs/config/*.go
	$(GO) fmt libs/file/*.go
	$(GO) fmt gcp.go

lint:
	$(GOLINT) libs/config
	$(GOLINT) libs/file
	$(GOLINT) gcp.go
