package main

import (
	"WeRide/user-service/db/postgres"
	"WeRide/user-service/internal/config"
	"WeRide/user-service/internal/repository"
	"WeRide/user-service/internal/service"
	"WeRide/user-service/pkg/logger"
	pb "WeRide/user-service/protoc/gen/go"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	logger2 "weride/internal/logger"
)

func main() {

	ctx := context.Background()

	ctx, err := logger.New(ctx)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to init config: %v", zap.Error(err))
	}

	pool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to init repository: %v", zap.Error(err))
	}
	ttlToken, err := cfg.GetAccessTokenTTL()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to ger access token", zap.Error(err))
	}
	repo := repository.NewRepository(pool, ttlToken, cfg.JwtSecret)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logger2.LoggerInterceptor))

	srv := service.New(repo)
	pb.RegisterAuthServer(grpcServer, srv)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", cfg.GRPCPort))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "Failed to listen", zap.Error(err))
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "gRPC server listening on "+cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
