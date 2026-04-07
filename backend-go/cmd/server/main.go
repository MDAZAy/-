package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"vpn-bot/backend-go/internal/config"
	"vpn-bot/backend-go/internal/db"
	"vpn-bot/backend-go/internal/handlers"
	"vpn-bot/backend-go/internal/middleware"
	"vpn-bot/backend-go/internal/providers"
	"vpn-bot/backend-go/internal/repositories"
	"vpn-bot/backend-go/internal/routes"
	"vpn-bot/backend-go/internal/services"
)

func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	database := db.Connect(cfg)
	db.Migrate(database)
	db.Seed(database, cfg)

	userRepo := repositories.NewUserRepository(database)
	planRepo := repositories.NewPlanRepository(database)
	subscriptionRepo := repositories.NewSubscriptionRepository(database)
	paymentRepo := repositories.NewPaymentRepository(database)
	vpnKeyRepo := repositories.NewVPNKeyRepository(database)

	paymentProvider := providers.NewPaymentProvider(cfg)
	vpnProvider := providers.NewVPNProvider(cfg)

	userService := services.NewUserService(userRepo)
	planService := services.NewPlanService(planRepo)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, planRepo, vpnKeyRepo, vpnProvider)
	paymentService := services.NewPaymentService(paymentRepo, planRepo, subscriptionService, paymentProvider)
	vpnService := services.NewVPNService(subscriptionRepo, vpnKeyRepo, vpnProvider)

	healthHandler := handlers.NewHealthHandler(cfg)
	userHandler := handlers.NewUserHandler(userService)
	planHandler := handlers.NewPlanHandler(planService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	vpnHandler := handlers.NewVPNHandler(vpnService)
	adminHandler := handlers.NewAdminHandler(userService, planService, subscriptionService, paymentService, vpnService)

	router := gin.New()
	router.Use(middleware.Logger(), middleware.Recovery())
	router.LoadHTMLGlob("internal/web/templates/*.tmpl")

	routes.Register(routes.Dependencies{
		Router:              router,
		Config:              cfg,
		HealthHandler:       healthHandler,
		UserHandler:         userHandler,
		PlanHandler:         planHandler,
		SubscriptionHandler: subscriptionHandler,
		PaymentHandler:      paymentHandler,
		VPNHandler:          vpnHandler,
		AdminHandler:        adminHandler,
	})

	go runExpirer(context.Background(), cfg.ExpirerInterval, subscriptionService)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend started on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func runExpirer(ctx context.Context, interval time.Duration, subscriptionService *services.SubscriptionService) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := subscriptionService.ExpireDue(ctx, time.Now()); err != nil {
				log.Printf("expirer job failed: %v", err)
			}
		}
	}
}
