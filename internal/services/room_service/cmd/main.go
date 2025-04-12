package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"we_ride/internal/services/room_service/config"
	"we_ride/internal/services/room_service/database"

	// api "weride/internal/services/room_service/api"

	"we_ride/internal/pkg/logger"

	"go.uber.org/zap"
	_ "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	ctx, _ = logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to read config:", zap.Error(err))
	}

	if err := database.RunMigrations(ctx, cfg.Postgres); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to run migrations:", zap.Error(err))
	}

	pool, err := database.New(ctx, cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to database:", zap.Error(err))
	}
	defer pool.Close()

	_, err = net.Listen("tcp", fmt.Sprintf("localhost:%s", cfg.GRPCPort))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to listen:", zap.Error(err))
	}
}
