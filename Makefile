GO ?= go
GOLINT ?= golint

build: fmt lint
	$(GO) build gcp.go

dep:
	$(GO) get
	$(GO) get -u github.com/golang/lint/golint

fmt:
	$(GO) fmt src/config/*.go
	$(GO) fmt src/file/*.go
	$(GO) fmt gcp.go

lint:
	$(GOLINT) src/config
	$(GOLINT) src/file
	$(GOLINT) gcp.go
