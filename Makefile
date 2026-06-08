-include .env
export

.PHONY: build docker-build run test

build:
	go build -a -o bin/api ./cmd/api

docker-build:
	docker build -t kodama:latest .

run:
	@go run ./cmd/api

test:
	go test -v -race ./...
