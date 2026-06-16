package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cristianriano/kodama/internal/config"
	"github.com/cristianriano/kodama/internal/infra/api"
	"github.com/cristianriano/kodama/internal/infra/logging"
	"github.com/cristianriano/kodama/internal/infra/telegram"
	"github.com/cristianriano/kodama/internal/usecase/chat"
)

const timeout = 5 * time.Second

func main() {
	logger := logging.Configure()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	httpClient := &http.Client{Timeout: cfg.TelegramAPITimeout}
	messenger := telegram.NewClient(cfg.TelegramBotToken, httpClient)
	chatUseCase := chat.NewUseCase(cfg.TelegramAllowedChatID, messenger)
	apiServer := api.NewServer(logger, cfg.TelegramWebhookSecret, chatUseCase, cfg.LogRequests)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           apiServer.Handler(),
		ReadHeaderTimeout: timeout,
		BaseContext:       func(net.Listener) context.Context { return ctx },
	}

	go func() {
		logger.Info("api server started", "addr", cfg.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdown(ctx, logger, httpServer, cfg.ShutdownTimeout)
}

func shutdown(ctx context.Context, logger *slog.Logger, server *http.Server, timeout time.Duration) {
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), timeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("api server shutdown failed", "error", err)
		return
	}

	logger.Info("api server stopped")
}
