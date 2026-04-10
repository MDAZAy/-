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
	VPNProviderUsername   string
	VPNProviderPassword   string
	VPNProviderInboundID  int
	VPNPublicHost         string
	VPNPublicPort         string
	VPNRealityServerName  string
	VPNRealityPublicKey   string
	VPNRealityShortID     string
	VPNFlow               string
	VPNFingerprint        string

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
		VPNProviderUsername:   getEnv("VPN_PROVIDER_USERNAME", ""),
		VPNProviderPassword:   getEnv("VPN_PROVIDER_PASSWORD", ""),
		VPNProviderInboundID:  getEnvInt("VPN_PROVIDER_INBOUND_ID", 0),
		VPNPublicHost:         getEnv("VPN_PROVIDER_PUBLIC_HOST", ""),
		VPNPublicPort:         getEnv("VPN_PROVIDER_PUBLIC_PORT", "443"),
		VPNRealityServerName:  getEnv("VPN_PROVIDER_REALITY_SERVER_NAME", ""),
		VPNRealityPublicKey:   getEnv("VPN_PROVIDER_REALITY_PUBLIC_KEY", ""),
		VPNRealityShortID:     getEnv("VPN_PROVIDER_REALITY_SHORT_ID", ""),
		VPNFlow:               getEnv("VPN_PROVIDER_FLOW", "xtls-rprx-vision"),
		VPNFingerprint:        getEnv("VPN_PROVIDER_FINGERPRINT", "chrome"),
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
