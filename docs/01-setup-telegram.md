# Telegram Setup

Official references:

- [Telegram Bots FAQ](https://core.telegram.org/bots/faq)
- [Telegram Bot API: setWebhook](https://core.telegram.org/bots/api#setwebhook)
- [Telegram Bot API: getWebhookInfo](https://core.telegram.org/bots/api#getwebhookinfo)

## 1. Create A Bot

Open [@BotFather](https://t.me/BotFather) in Telegram and run:

```text
/newbot
Hi
```

BotFather will ask for:

- A display name, for example `Kodama`
- A username ending in `bot`, for example `my_kodama_bot`

BotFather will return a token. Keep it secret.

## 2. Get Your Chat ID

Send your bot a message in Telegram, for example:

```text
/start
```

Then run:

```sh
export TELEGRAM_BOT_TOKEN="replace-with-botfather-token"

curl "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getUpdates"
```

Find the chat ID in the response:

```json
{
  "message": {
    "chat": {
      "id": 123456789
    }
  }
}
```

Use that value in `.env`:

```sh
TELEGRAM_ALLOWED_CHAT_ID=123456789
```

## 3. Choose A Webhook Secret

Telegram can send a secret header with every webhook request. Kodama validates that header before processing the update.

Generate a secret:

```sh
openssl rand -hex 32
```

Use that value in `.env`:

```sh
TELEGRAM_WEBHOOK_SECRET=replace-with-generated-secret
```

Telegram's `secret_token` only allows letters, numbers, `_`, and `-`, so hex is a simple safe format.

## 4. Expose Kodama

Telegram webhooks require a public HTTPS URL. A typical home-server setup is:

```text
Telegram -> https://your-domain.com/telegram/webhook -> reverse proxy -> Kodama :8080
```

## 5. Register The Webhook

Webhook registration is a deployment step. Kodama does not register the webhook on startup.

```sh
export TELEGRAM_BOT_TOKEN="replace-with-botfather-token"
export TELEGRAM_WEBHOOK_SECRET="replace-with-generated-secret"
export KODAMA_PUBLIC_URL="https://your-domain.com"

curl -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setWebhook" \
  -d "url=${KODAMA_PUBLIC_URL}/telegram/webhook" \
  -d "secret_token=${TELEGRAM_WEBHOOK_SECRET}" \
  -d 'allowed_updates=["message"]' \
  -d "drop_pending_updates=true"
```

## 6. Verify

```sh
curl "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getWebhookInfo"
```

Check that the configured URL is correct and that Telegram reports no recent delivery errors.
