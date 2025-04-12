package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"we_ride/internal/pkg/logger"
	"we_ride/internal/services/room_service/config"
	room_service_transport "we_ride/internal/services/room_service/transport"

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
	/*
		if err := database.RunMigrations(ctx, cfg.Postgres); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to run migrations:", zap.Error(err))
		}

		pool, err := database.New(ctx, cfg.Postgres)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to connect to database:", zap.Error(err))
		}
		defer pool.Close()
	*/

	roomS := room_service_transport.New()
	r := room_service_transport.NewRouter(ctx, cfg, *roomS)

	r.Run(ctx)

	select {
	case <-ctx.Done():
		r.Server.GracefulStop()
		// pool.Close()
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Server Stopped")
	}
}
