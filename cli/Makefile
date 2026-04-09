.PHONY: build test lint

build:
	go build -o bin/securedbx ./cmd/securedbx

test:
	go test ./... -v

lint:
	golangci-lint run
