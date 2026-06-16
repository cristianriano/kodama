package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr                  string
	LogRequests           bool
	TelegramAllowedChatID int64
	TelegramWebhookSecret string
	ShutdownTimeout       time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Addr:                  env("ADDR", ":8080"),
		LogRequests:           false,
		TelegramWebhookSecret: env("TELEGRAM_WEBHOOK_SECRET", ""),
		ShutdownTimeout:       10 * time.Second,
	}

	logRequests := env("LOG_REQUESTS", "")
	if logRequests != "" {
		parsedLogRequests, err := strconv.ParseBool(logRequests)
		if err != nil {
			return Config{}, fmt.Errorf("parse LOG_REQUESTS: %w", err)
		}
		cfg.LogRequests = parsedLogRequests
	}

	allowedChatID := env("TELEGRAM_ALLOWED_CHAT_ID", "")
	if allowedChatID == "" {
		return Config{}, errors.New("TELEGRAM_ALLOWED_CHAT_ID is required")
	}

	parsedChatID, err := strconv.ParseInt(allowedChatID, 10, 64)
	if err != nil {
		return Config{}, fmt.Errorf("parse TELEGRAM_ALLOWED_CHAT_ID: %w", err)
	}
	cfg.TelegramAllowedChatID = parsedChatID

	if cfg.TelegramWebhookSecret == "" {
		return Config{}, errors.New("TELEGRAM_WEBHOOK_SECRET is required")
	}

	return cfg, nil
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
