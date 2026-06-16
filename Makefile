-include .env
export

.PHONY: build docker-build docker-deploy get-webhook run send set-webhook test

build:
	go build -a -o bin/api ./cmd/api

docker-build:
	docker build -t kodama:latest .

docker-deploy:
	-docker-compose down
	docker-compose build
	docker-compose up -d

run:
	@go run ./cmd/api

test:
	go test -v -race ./...

set-webhook:
	@test -n "${TELEGRAM_BOT_TOKEN}" || (echo "TELEGRAM_BOT_TOKEN is required" && exit 1)
	@test -n "${KODAMA_PUBLIC_URL}" || (echo "KODAMA_PUBLIC_URL is required" && exit 1)
	@test -n "${TELEGRAM_WEBHOOK_SECRET}" || (echo "TELEGRAM_WEBHOOK_SECRET is required" && exit 1)
	@curl -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setWebhook" \
		-d "url=${KODAMA_PUBLIC_URL}/telegram/webhook" \
		-d "secret_token=${TELEGRAM_WEBHOOK_SECRET}" \
		-d 'allowed_updates=["message"]' \
		-d "drop_pending_updates=true"

get-webhook:
	@test -n "${TELEGRAM_BOT_TOKEN}" || (echo "TELEGRAM_BOT_TOKEN is required" && exit 1)
	@curl "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getWebhookInfo"

MSG ?= Hi
send:
	@curl -X POST -H "X-Telegram-Bot-Api-Secret-Token: ${TELEGRAM_WEBHOOK_SECRET}" localhost:8080/telegram/webhook \
		-d '{"message": {"text": "${MSG}", "chat": {"id": ${TELEGRAM_ALLOWED_CHAT_ID} } } }'
