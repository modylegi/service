package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `env:"ENV" env-required:"true"`
	UseCache   bool   `env:"USE_CACHE" env-required:"true"`
	HttpServer HttpServer
	DB         DB
	Cache      Cache
}

type HttpServer struct {
	Port         int           `env:"SERVER_PORT" env-required:"true"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" env-required:"true"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" env-required:"true"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" env-required:"true"`
}

type DB struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-required:"true"`
	Database string `env:"DB_DATABASE" env-required:"true"`
	Username string `env:"DB_USERNAME" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

type Cache struct {
	Host     string `env:"CACHE_HOST" env-required:"true"`
	Port     int    `env:"CACHE_PORT" env-required:"true"`
	Password string `env:"CACHE_PASSWORD" env-required:"true"`
}

func Load() (*Config, error) {
	// configPath := fetchConfigPath()
	// if configPath == "" {
	// 	return nil, errors.New("config path is empty")
	// }

	// if _, err := os.Stat(configPath); os.IsNotExist(err) {
	// 	return nil, fmt.Errorf("config file does not exist: %w", err)
	// }

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}
	// if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
	// 	return nil, fmt.Errorf("cannot read config: %w", err)
	// }
	return &cfg, nil
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
