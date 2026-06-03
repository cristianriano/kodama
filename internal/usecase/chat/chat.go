package chat

import (
	"context"
	"strings"
)

type Messenger interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type UseCase struct {
	allowedChatID int64
	messenger     Messenger
}

func NewUseCase(allowedChatID int64, messenger Messenger) *UseCase {
	return &UseCase{
		allowedChatID: allowedChatID,
		messenger:     messenger,
	}
}

type Request struct {
	ChatID int64
	Text   string
}

func (u *UseCase) ReceiveMessage(ctx context.Context, req Request) error {
	if u.allowedChatID != 0 && req.ChatID != u.allowedChatID {
		return nil
	}

	text := strings.TrimSpace(req.Text)
	if text == "" {
		return nil
	}

	if err := u.messenger.SendMessage(ctx, req.ChatID, text); err != nil {
		return err
	}

	return nil
}
