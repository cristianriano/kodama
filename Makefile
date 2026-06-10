-include .env
export

.PHONY: build docker-build run send test

build:
	go build -a -o bin/api ./cmd/api

docker-build:
	docker build -t kodama:latest .

run:
	@go run ./cmd/api

test:
	go test -v -race ./...

MSG ?= Hi
send:
	@curl -X POST -H "X-Telegram-Bot-Api-Secret-Token: ${TELEGRAM_WEBHOOK_SECRET}" localhost:8080/telegram/webhook \
		-d '{"message": {"text": "${MSG}", "chat": {"id": ${TELEGRAM_ALLOWED_CHAT_ID} } } }'