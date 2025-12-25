package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServiceName string
	Port        int

	RedisURL string

	ShutdownTimeoutSeconds int
}

func Load() *Config {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		log.Fatal("invalid PORT")
	}

	shutdownTimeout, err := strconv.Atoi(getEnv("SHUTDOWN_TIMEOUT_SECONDS", "10"))
	if err != nil {
		log.Fatal("invalid SHUTDOWN_TIMEOUT_SECONDS")
	}

	cfg := &Config{
		ServiceName: "inventory",
		Port:        port,
		RedisURL:   mustGetEnv("REDIS_URL"),
		ShutdownTimeoutSeconds: shutdownTimeout,
	}

	return cfg
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
