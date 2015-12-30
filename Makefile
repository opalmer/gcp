GO ?= go
GOLINT ?= golint

build: fmt lint
	$(GO) build gcp.go

dep:
	$(GO) get
	$(GO) get -u github.com/golang/lint/golint

fmt:
	$(GO) fmt gcp.go

lint:
	$(GOLINT) gcp.go
