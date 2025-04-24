package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"

	"we_ride/internal/services/payment_service/db/postgres"
	"we_ride/internal/services/payment_service/internal/service"
	pb "we_ride/internal/services/payment_service/protoc/gen"
)

type Config struct {
	// GRPC    config.GRPCConfig `yaml:"grpc"`
	Postgres postgres.Config  `yaml:"postgres"`
	// UsersGRPC config.GRPCConfig `yaml:"users_grpc"`
}

func main() {
	var cfg Config
	if err := cleanenv.ReadConfig("./config/local.yaml", &cfg); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	
	// Database connection
	pool, err := postgres.New(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatal("Postgres init error", err)
	}
	defer pool.Close()

	// UserService connection
	// userConn, err := grpc.Dial(
	// 	fmt.Sprintf("%s:%s", cfg.UsersGRPC.Host, cfg.UsersGRPC.Port),
	// 	grpc.WithInsecure(),
	// )
	// if err != nil {
	// 	log.Fatal("UserService connection failed", zap.Error(err))
	// }
	// defer userConn.Close()

	// GRPC Server
	srv := grpc.NewServer()
	paymentService, err := service.New(pool)
	pb.RegisterPaymentServiceServer(srv, paymentService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", "50053"))
	if err != nil {
		log.Fatal("Listener error", err)
	}

	log.Print("Starting PaymentService")
	if err := srv.Serve(listener); err != nil {
		log.Fatal("Server error", err)
	}
}