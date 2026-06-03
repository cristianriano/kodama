# Kodama

Personal AI assistant built incrementally in Go.

## Current slice

Kodama currently runs a small HTTP API for Telegram webhooks. When it receives a text message from the configured chat ID, it echoes that message to the console through a temporary messenger implementation.

Telegram webhook registration must be done manually as part of deployment.

## Run

```sh
export TELEGRAM_ALLOWED_CHAT_ID="123456789"
export TELEGRAM_WEBHOOK_SECRET="shared-secret"

make run
```

Optional:

```sh
export ADDR=":8080"
```

Webhook endpoint:

```text
POST /telegram/webhook
```
