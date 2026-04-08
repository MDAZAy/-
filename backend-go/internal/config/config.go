package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	Port            string
	DatabaseURL     string
	PublicBaseURL   string
	AdminToken      string
	PaymentProvider string
	VPNProvider     string

	CloudPaymentsPublicID string
	CloudPaymentsAPIToken string
	VPNProviderURL        string
	VPNProviderToken      string

	SeedPlansPath   string
	ExpirerInterval time.Duration
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("config: .env not loaded: %v", err)
	}

	return Config{
		AppEnv:                getEnv("APP_ENV", "development"),
		Port:                  getEnv("PORT", "8080"),
		DatabaseURL:           getEnv("DATABASE_URL", "data/app.db"),
		PublicBaseURL:         getEnv("PUBLIC_BASE_URL", "http://localhost:8080"),
		AdminToken:            getEnv("ADMIN_TOKEN", "change-me-admin-token"),
		PaymentProvider:       getEnv("PAYMENT_PROVIDER", "mock"),
		VPNProvider:           getEnv("VPN_PROVIDER", "mock"),
		CloudPaymentsPublicID: getEnv("CLOUDPAYMENTS_PUBLIC_ID", ""),
		CloudPaymentsAPIToken: getEnv("CLOUDPAYMENTS_API_SECRET", ""),
		VPNProviderURL:        getEnv("VPN_PROVIDER_ENDPOINT", ""),
		VPNProviderToken:      getEnv("VPN_PROVIDER_TOKEN", ""),
		SeedPlansPath:         getEnv("SEED_PLANS_PATH", "seed/plans.json"),
		ExpirerInterval:       time.Duration(getEnvInt("EXPIRER_INTERVAL_SECONDS", 60)) * time.Second,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
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
