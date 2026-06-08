# Kodama

Personal AI assistant built incrementally in Go.

## Overview

Kodama is designed as a small Go service for a personal AI assistant. It receives Telegram webhook updates, routes them through explicit use cases, and sends responses through infrastructure adapters.

## Local Setup

1. Create a local `.env` file and fill the required values:

```sh
cp .env.example .env
```

2. Generate a webhook secret with:

```sh
openssl rand -hex 32
```

Telegram's webhook `secret_token` only allows letters, numbers, `_`, and `-`, so hex is a simple safe format.

`make run` loads `.env` through the Makefile. If you run the binary directly, provide the variables through the shell, systemd, Docker, or your process manager.

3. Send a message to confirm everything works fine

```sh
curl -X POST -H "X-Telegram-Bot-Api-Secret-Token: <TELEGRAM_WEBHOOK_SECRET>" localhost:8080/telegram/webhook \
-d '{"message": {"text": "Hi", "chat": {"id": TELEGRAM_ALLOWED_CHAT_ID}}}'
```

## Development

```sh
make run
make test
make build
make docker-build
```

Webhook endpoint:

```text
POST /telegram/webhook
```
