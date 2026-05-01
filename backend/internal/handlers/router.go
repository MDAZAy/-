package handlers

import (
	"net/http"
	"time"

	"autoservice/backend/internal/auth"
	"autoservice/backend/internal/config"
	"autoservice/backend/internal/dto"
	"autoservice/backend/internal/middleware"
	"autoservice/backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	authService        *services.AuthService
	catalogService     *services.CatalogService
	appointmentService *services.AppointmentService
}

func NewRouter(cfg config.Config, jwtManager *auth.Manager, authService *services.AuthService, catalogService *services.CatalogService, appointmentService *services.AppointmentService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Idempotency-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middleware.NewRateLimit(cfg.RateLimitRequests, cfg.RateLimitWindow))

	handler := &HTTPHandler{
		authService:        authService,
		catalogService:     catalogService,
		appointmentService: appointmentService,
	}

	router.GET("/health", func(c *gin.Context) {
		success(c, gin.H{"status": "ok"})
	})
	router.StaticFile("/swagger/swagger.yaml", "docs/swagger.yaml")

	api := router.Group("/api/v1")

	authGroup := api.Group("/auth")
	authGroup.Use(middleware.NewRateLimit(cfg.AuthRateLimitRequests, cfg.AuthRateLimitWindow))
	authGroup.POST("/register", handler.register)
	authGroup.POST("/login", handler.login)
	authGroup.POST("/refresh", handler.refresh)
	authGroup.POST("/logout", handler.logout)

	api.GET("/services", handler.listServices)
	api.GET("/service-categories", handler.listCategories)
	api.GET("/appointments/available-slots", handler.availableSlots)

	protected := api.Group("/")
	protected.Use(middleware.AuthRequired(jwtManager))
	protected.GET("/me", handler.me)
	protected.GET("/vehicles/my", handler.listVehicles)
	protected.POST("/vehicles", handler.createVehicle)
	protected.POST("/appointments", handler.createAppointment)
	protected.GET("/appointments/my", handler.listMyAppointments)

	admin := protected.Group("/")
	admin.Use(middleware.RequireRole("admin"))
	admin.GET("/appointments", handler.listAllAppointments)
	admin.GET("/admin/dashboard", handler.dashboard)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, dto.Envelope{Success: false, Error: "route not found", Code: "not_found"})
	})

	return router
}

func success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, dto.Envelope{Success: true, Data: data})
}

func failure(c *gin.Context, appErr *services.AppError) {
	c.JSON(appErr.Status, dto.Envelope{Success: false, Error: appErr.Message, Code: appErr.Code})
}
