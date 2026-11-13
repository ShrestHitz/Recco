package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	Env            string
	Port           string
	DatabaseURL    string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	DBMaxConns     int32
	DBMinConns     int32
	DBConnLifetime time.Duration
}

// Load loads application configuration with production-safe defaults.
func Load() *Config {
	return &Config{
		Env:            getEnv("APP_ENV", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/recco?sslmode=disable"),
		ReadTimeout:    getEnvDuration("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:   getEnvDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:    getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		DBMaxConns:     int32(getEnvInt("DB_MAX_CONNS", 25)),
		DBMinConns:     int32(getEnvInt("DB_MIN_CONNS", 5)),
		DBConnLifetime: getEnvDuration("DB_CONN_LIFETIME", 1*time.Hour),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := time.ParseDuration(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}
