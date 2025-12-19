.PHONY: build test lint clean install

VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/gham ./cmd/gham

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

install:
	go install $(LDFLAGS) ./cmd/gham
