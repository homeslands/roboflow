package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/tuanvumaihuynh/roboflow/internal/server"
	"github.com/tuanvumaihuynh/roboflow/pkg/config"
	"github.com/tuanvumaihuynh/roboflow/pkg/logs"
)

func main() {
	// Setup UTC time
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	// Load configuration
	cfg := config.MustLoadConfig()

	// Create logger
	logger := logs.NewZapLogger(cfg.IsProd())
	defer logger.Sync()

	// Connect to database
	conn, err := pgxpool.New(context.Background(), cfg.PostgresDsn)
	if err != nil {
		logger.Fatal("Error connecting to db", zap.Error(err))
	}
	defer conn.Close()

	if err := conn.Ping(context.Background()); err != nil {
		logger.Fatal("Error pinging to db", zap.Error(err))
	}

	// Run server
	cleanupFn := server.Run(cfg, conn, logger)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cleanupFn(ctx)

	logger.Info("Server stopped")
}
