package config

import (
	"WeRide/user-service/db/postgres"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Postgres postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`

	GRPCPort string `yaml:"GRPC_PORT" env:"GRPC_PORT"  envDefault:"50051"`
	RESRPort string `yaml:"REST_PORT" env:"REST_PORT" envDefault:"8081"`

	JWTAccessTokenTTL string `yaml:"JWT_ACCESS_TOKEN_TTL" env:"JWT_ACCESS_TOKEN_TTL" envDefault:"15m"`
	JwtSecret         string `yaml:"JWT_SECRET" env:"JWT_SECRET" envDefault:"secret"`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("./config/local.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) GetAccessTokenTTL() (time.Duration, error) {
	return time.ParseDuration(c.JWTAccessTokenTTL)
}
