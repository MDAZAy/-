package main

import (
	"log"

	"autoservice/backend/internal/auth"
	"autoservice/backend/internal/config"
	"autoservice/backend/internal/database"
	"autoservice/backend/internal/handlers"
	"autoservice/backend/internal/repositories"
	"autoservice/backend/internal/services"
)

func main() {
	cfg := config.Load()

	db, sqlDB, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	if cfg.AutoMigrate {
		if err := database.ApplyMigrations(sqlDB, cfg.MigrationsDir); err != nil {
			log.Fatalf("migrations failed: %v", err)
		}
	}

	jwtManager := auth.NewManager(cfg)
	if err := database.Seed(db); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	repo := repositories.New(db)
	authService := services.NewAuthService(repo, jwtManager, cfg)
	catalogService := services.NewCatalogService(repo, cfg)
	appointmentService := services.NewAppointmentService(repo, cfg)

	router := handlers.NewRouter(cfg, jwtManager, authService, catalogService, appointmentService)

	log.Printf("autoservice API listening on :%s", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
