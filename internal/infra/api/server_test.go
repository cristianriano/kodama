package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cristianriano/kodama/internal/usecase/chat"
)

type fakeMessenger struct {
	chatID int64
	text   string
	calls  int
}

func (f *fakeMessenger) SendMessage(_ context.Context, chatID int64, text string) error {
	f.chatID = chatID
	f.text = text
	f.calls++
	return nil
}

func TestTelegramWebhook(t *testing.T) {
	tests := []struct {
		name          string
		headerSecret  string
		body          string
		wantStatus    int
		wantSendCalls int
		wantChatID    int64
		wantText      string
	}{
		{
			name:         "valid request sends same message and returns 200",
			headerSecret: "secret",
			body: `{
				"message": {
					"text": "hello",
					"chat": {"id": 123}
				}
			}`,
			wantStatus:    http.StatusOK,
			wantSendCalls: 1,
			wantChatID:    123,
			wantText:      "hello",
		},
		{
			name:          "invalid header returns 401",
			headerSecret:  "wrong",
			body:          `{}`,
			wantStatus:    http.StatusUnauthorized,
			wantSendCalls: 0,
		},
		{
			name:         "request from another chat returns 200",
			headerSecret: "secret",
			body: `{
				"message": {
					"text": "hello",
					"chat": {"id": 456}
				}
			}`,
			wantStatus:    http.StatusOK,
			wantSendCalls: 0,
		},
		{
			name:          "ignore non message update",
			headerSecret:  "secret",
			body:          `{"update_id": 1}`,
			wantStatus:    http.StatusOK,
			wantSendCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messenger := &fakeMessenger{}
			chatUseCase := chat.NewUseCase(123, messenger)
			server := NewServer(testLogger(), "secret", chatUseCase, false)

			req := httptest.NewRequest(http.MethodPost, "/telegram/webhook", strings.NewReader(tt.body))
			req.Header.Set(telegramSecretHeader, tt.headerSecret)
			rec := httptest.NewRecorder()

			server.Handler().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if messenger.calls != tt.wantSendCalls {
				t.Fatalf("SendMessage calls = %d, want %d", messenger.calls, tt.wantSendCalls)
			}
			if tt.wantSendCalls == 0 {
				return
			}
			if messenger.chatID != tt.wantChatID {
				t.Fatalf("chatID = %d, want %d", messenger.chatID, tt.wantChatID)
			}
			if messenger.text != tt.wantText {
				t.Fatalf("text = %q, want %q", messenger.text, tt.wantText)
			}
		})
	}
}

func TestHealthz(t *testing.T) {
	server := NewServer(testLogger(), "secret", nil, true)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "OK\n" {
		t.Fatalf("body = %q, want %q", rec.Body.String(), "ok\n")
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
