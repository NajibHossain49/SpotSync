package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all environment-driven settings.
type Config struct {
	DatabaseURL    string
	JWTSecret      string
	JWTExpiryHours int
	Port           string
}

// Load reads the .env file (if present) and environment variables.
func Load() *Config {
	// In production (Render/Railway) env vars are injected directly,
	// so a missing .env file is not fatal.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment variables")
	}

	expiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		expiry = 24
	}

	cfg := &Config{
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTExpiryHours: expiry,
		Port:           getEnv("PORT", "8080"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
