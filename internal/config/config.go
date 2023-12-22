package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"prod"`
	StoragePath string        `yaml:"storage_path" env-reauired:"true"`
	Token_ttl   time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPS        GRPS_conf     `yaml:"grpc"`
}

type GRPS_conf struct {
	Port    int    `yaml:"port" env-drfault:"8000"`
	TimeOut string `yaml:"timeout" env-default:"4s"`
}

func MustLoad() Config {
	// TODO: переделать
	configPath := "D:/go_path/src/jwt_auth_gRPC/sso/config/local.yaml"
	err := os.Setenv("CONFIG_PATH", configPath)
	if err != nil {
		log.Fatal("Failed to set CONFIG_PATH")
	}

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG PATH is not set")
	}

	// checl if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
