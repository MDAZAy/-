package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"lua-agent/backend/internal/agent"
	"lua-agent/backend/internal/api"
	"lua-agent/backend/internal/config"
	"lua-agent/backend/internal/llm"
	"lua-agent/backend/internal/storage"
	"lua-agent/backend/internal/validator"
	"lua-agent/backend/pkg/logger"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config/config.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Logging.Level)
	ctx := context.Background()

	repo, err := storage.NewPostgresRepository(ctx, storage.PostgresConfig{
		DSN: cfg.Storage.DSN,
	})
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	llmClient := llm.NewClient(cfg.LLM.Endpoint, cfg.LLM.Model, cfg.LLM.Timeout)
	_ = llm.NewEmbedClient(cfg.Embedding.Endpoint, cfg.Embedding.Timeout)
	validate := validator.New(cfg.Validator.Timeout)

	service := agent.NewService(llmClient, repo, validate, agent.LLMConfig{
		Model:      cfg.LLM.Model,
		NumCtx:     cfg.LLM.NumCtx,
		NumPredict: cfg.LLM.NumPredict,
		Batch:      cfg.LLM.Batch,
		Parallel:   cfg.LLM.Parallel,
	})

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      api.NewRouter(service, repo, log),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("agent backend started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped with error", "error", err)
			stop()
		}
	}()

	<-stopCtx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
	}
}
