.PHONY: build test lint fmt vet clean install install-git-hooks tidy site-build site-validate

BINARY_NAME=scout
VERSION?=$(shell git describe --tags --abbrev=0 --match 'v[0-9]*.[0-9]*.[0-9]*' 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-X github.com/jamesonstone/scout/pkg/cli.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/scout

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME).exe ./cmd/scout

install:
	go install $(LDFLAGS) ./cmd/scout

install-git-hooks:
	@if [ -f .githooks/pre-commit ]; then chmod +x .githooks/pre-commit && git config core.hooksPath .githooks; else echo "no .githooks/pre-commit present"; fi

test:
	go test -v ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin/
	go clean

tidy:
	go mod tidy

site-build:
	go run ./cmd/scout site build --data-dir . --out-dir public --base-path /scout/

site-validate:
	go run ./cmd/scout site validate --out-dir public --base-path /scout/

all: fmt vet test build
