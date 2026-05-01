package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                string
	HTTPPort              string
	MySQLDSN              string
	AccessSecret          string
	RefreshSecret         string
	AccessTTL             time.Duration
	RefreshTTL            time.Duration
	AutoMigrate           bool
	MigrationsDir         string
	DisplayTimezone       string
	RateLimitRequests     int
	RateLimitWindow       time.Duration
	AuthRateLimitRequests int
	AuthRateLimitWindow   time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:                getEnv("APP_ENV", "development"),
		HTTPPort:              getEnv("HTTP_PORT", "8080"),
		MySQLDSN:              getEnv("MYSQL_DSN", "autoservice:autoservice@tcp(localhost:3306)/autoservice?charset=utf8mb4&parseTime=True&loc=UTC"),
		AccessSecret:          getEnv("JWT_ACCESS_SECRET", "change-me-access-secret"),
		RefreshSecret:         getEnv("JWT_REFRESH_SECRET", "change-me-refresh-secret"),
		AccessTTL:             time.Duration(getInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
		RefreshTTL:            time.Duration(getInt("JWT_REFRESH_TTL_HOURS", 720)) * time.Hour,
		AutoMigrate:           getBool("AUTO_MIGRATE", true),
		MigrationsDir:         getEnv("MIGRATIONS_DIR", "migrations"),
		DisplayTimezone:       getEnv("DISPLAY_TIMEZONE", "Europe/Helsinki"),
		RateLimitRequests:     getInt("RATE_LIMIT_REQUESTS", 120),
		RateLimitWindow:       time.Duration(getInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
		AuthRateLimitRequests: getInt("AUTH_RATE_LIMIT_REQUESTS", 10),
		AuthRateLimitWindow:   time.Duration(getInt("AUTH_RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
	}

	if cfg.MySQLDSN == "" {
		log.Fatal("MYSQL_DSN is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
