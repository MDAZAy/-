package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"vpn-bot/backend-go/internal/config"
)

func Connect(cfg config.Config) *gorm.DB {
	databaseDir := filepath.Dir(cfg.DatabaseURL)
	if databaseDir != "." {
		if err := os.MkdirAll(databaseDir, 0o755); err != nil {
			log.Fatalf("db mkdir failed: %v", err)
		}
	}

	database, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	return database
}
