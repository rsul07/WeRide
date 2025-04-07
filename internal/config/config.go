package config

import (
	"fmt"
	"weride/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB postgres.DBConfig `env:"POSTGRES" env-default:"POSTGRES" yaml:"POSTGRES"`

	GRPCPort string `env:"GRPC_PORT" env-default:"50051"   yaml:"GRPC_PORT"`
	GRPCHost string `env:"GRPC_HOST" env-default:"0.0.0.0" yaml:"GRPC_HOST"`

	RESTPort string `env:"REST_PORT" env-default:"8081"    yaml:"REST_PORT"`
	RESTHost string `env:"REST_HOST" env-default:"0.0.0.0" yaml:"REST_HOST"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("./config/config.yaml", &cfg)
	if err != nil {
		return nil, fmt.Errorf("config.Newconfig: %w", err)
	}

	return &cfg, nil
}
