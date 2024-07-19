package app

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env                string
	ServerPort         string
	ServerIdleTimeout  time.Duration
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	DBDatabase         string
	DBPassword         string
	DBUsername         string
	DBHost             string
	DBPort             string
	Cache              bool
	CacheHost          string
	CachePort          string
	CachePassword      string
}

func parseEnvDuration(envVar string, defaultVal time.Duration) time.Duration {
	if val, ok := os.LookupEnv(envVar); ok {
		if parsedVal, err := time.ParseDuration(val); err == nil {
			return parsedVal
		}
	}
	return defaultVal
}

func parseEnvInt(envVar string, defaultVal int) int {
	if val, ok := os.LookupEnv(envVar); ok {
		if parsedVal, err := strconv.Atoi(val); err == nil {
			return parsedVal
		}
	}
	return defaultVal
}

func getConfig() *Config {
	return &Config{
		Env:                os.Getenv("ENV"),
		ServerPort:         os.Getenv("SERVER_PORT"),
		ServerIdleTimeout:  parseEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		ServerReadTimeout:  parseEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		ServerWriteTimeout: parseEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		DBDatabase:         os.Getenv("DB_DATABASE"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBUsername:         os.Getenv("DB_USERNAME"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		Cache:              os.Getenv("CACHE") == "true",
		CacheHost:          os.Getenv("CACHE_HOST"),
		CachePort:          os.Getenv("CACHE_PORT"),
		CachePassword:      os.Getenv("CACHE_PASSWORD"),
	}
}
