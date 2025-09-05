package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/RoGogDBD/kasp/internal/config"
	"github.com/RoGogDBD/kasp/internal/handlers"
	"github.com/RoGogDBD/kasp/internal/logger"
	"github.com/RoGogDBD/kasp/internal/repository"
	"github.com/RoGogDBD/kasp/internal/service"
)

func main() {
	config.ParseFlags()
	if err := run(); err != nil {
		logger.Error(fmt.Sprintf("server error: %v", err))
	}
}

func run() error {
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf(
		"config: run_addr=%s workers=%d queue_size=%d log_level=%s",
		config.FlagRunAddr, config.FlagWorkers, config.FlagQueueSize, config.FlagLogLevel,
	))

	queue := repository.NewQueue(config.FlagQueueSize)
	storage := repository.NewStorage()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go service.StartWorkers(ctx, config.FlagWorkers, queue, storage)

	mux := http.NewServeMux()
	mux.Handle("/enqueue", logger.RequestLogger(handlers.EnqueueHandler(queue, storage)))
	mux.HandleFunc("/healthz", handlers.HealthHandler())

	srv := &http.Server{
		Addr:    config.FlagRunAddr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("listen error: %v", err))
		}
	}()

	logger.Info(fmt.Sprintf("Running server on address %s", config.FlagRunAddr))

	<-ctx.Done()
	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("Server forced to shutdown: %w", err)
	}

	return nil
}
