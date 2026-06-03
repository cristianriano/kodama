package telegram

import (
	"context"
	"log/slog"
)

type ConsoleMessenger struct {
	logger *slog.Logger
}

func NewConsoleMessenger(logger *slog.Logger) *ConsoleMessenger {
	return &ConsoleMessenger{logger: logger}
}

func (m *ConsoleMessenger) SendMessage(ctx context.Context, chatID int64, text string) error {
	m.logger.InfoContext(ctx, "telegram send message", "chat_id", chatID, "text", text)
	return nil
}
