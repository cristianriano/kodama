package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cristianriano/kodama/internal/usecase/chat"
)

const telegramSecretHeader = "X-Telegram-Bot-Api-Secret-Token"

type ChatHandler interface {
	ReceiveMessage(ctx context.Context, req chat.Request) error
}

type Server struct {
	logger         *slog.Logger
	logRequests    bool
	telegramSecret string
	chatHandler    ChatHandler
}

func NewServer(logger *slog.Logger, telegramSecret string, chatHandler ChatHandler, logRequests bool) *Server {
	return &Server{
		logger:         logger,
		logRequests:    logRequests,
		telegramSecret: telegramSecret,
		chatHandler:    chatHandler,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("POST /telegram/webhook", s.handleTelegramWebhook)
	if s.logRequests {
		return requestLogger(s.logger, mux)
	}
	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK\n"))
}

func (s *Server) handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	if s.telegramSecret != "" && r.Header.Get(telegramSecretHeader) != s.telegramSecret {
		http.Error(w, "invalid telegram secret", http.StatusUnauthorized)
		return
	}

	var update telegramUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req, ok := update.ChatRequest()
	if !ok {
		// Acknowledge update types we do not handle so Telegram does not retry them.
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := s.chatHandler.ReceiveMessage(r.Context(), req); err != nil {
		s.logger.ErrorContext(r.Context(), "handle telegram webhook", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type telegramUpdate struct {
	Message *telegramMessage `json:"message"`
}

func (u telegramUpdate) ChatRequest() (chat.Request, bool) {
	if u.Message == nil {
		return chat.Request{}, false
	}
	if u.Message.Chat.ID == 0 || u.Message.Text == "" {
		return chat.Request{}, false
	}

	return chat.Request{
		ChatID: u.Message.Chat.ID,
		Text:   u.Message.Text,
	}, true
}

type telegramMessage struct {
	Text string       `json:"text"`
	Chat telegramChat `json:"chat"`
}

type telegramChat struct {
	ID int64 `json:"id"`
}
