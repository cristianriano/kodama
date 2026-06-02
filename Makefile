.PHONY: build run test

build:
	go build -a -o bin/api ./cmd/api

run:
	@go run ./cmd/api

test:
	go test -v -race ./...
