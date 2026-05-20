package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	AdminAPIURL string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", getEnv("DB_USER", "postgres"), getEnv("DB_PASSWORD", "postgres"), getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432"), getEnv("DB_NAME", "calendar"))),
		AdminAPIURL: getEnv("ADMIN_API_URL", "http://localhost:8081"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
