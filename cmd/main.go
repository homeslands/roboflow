package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tuanvumaihuynh/roboflow/internal/application"
	"github.com/tuanvumaihuynh/roboflow/pkg/config"
	"github.com/tuanvumaihuynh/roboflow/pkg/logs"
)

func main() {
	// Context
	ctx := context.Background()

	// Setup UTC time
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	// Load configuration
	cfg := config.MustLoadConfig()

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}
	handler := logs.Handler{Handler: slog.NewTextHandler(os.Stdout, opts)}
	logger := slog.New(handler)

	// Connect to database
	connCfg, err := pgxpool.ParseConfig(cfg.PostgresDsn)
	if err != nil {
		panic(err)
	}
	conn, err := pgxpool.NewWithConfig(ctx, connCfg)
	if err != nil {
		logger.Error("Error connecting to db", slog.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		logger.Error("Error pinging to db", slog.Any("error", err))
		os.Exit(1)
	}

	// Run server
	cleanupFn := application.Run(*cfg, conn, logger)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-stop

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	cleanupFn(ctx)

	logger.Info("Server stopped")
}
