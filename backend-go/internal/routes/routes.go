package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vpn-bot/backend-go/internal/config"
	"vpn-bot/backend-go/internal/handlers"
	"vpn-bot/backend-go/internal/middleware"
)

type Dependencies struct {
	Router              *gin.Engine
	Config              config.Config
	HealthHandler       *handlers.HealthHandler
	UserHandler         *handlers.UserHandler
	PlanHandler         *handlers.PlanHandler
	SubscriptionHandler *handlers.SubscriptionHandler
	PaymentHandler      *handlers.PaymentHandler
	VPNHandler          *handlers.VPNHandler
	AdminHandler        *handlers.AdminHandler
}

func Register(dep Dependencies) {
	dep.Router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/health")
	})
	dep.Router.GET("/health", dep.HealthHandler.Get)

	api := dep.Router.Group("/api/v1")
	{
		api.POST("/users/ensure", dep.UserHandler.Ensure)
		api.GET("/plans", dep.PlanHandler.List)
		api.POST("/plans", middleware.AdminAuth(dep.Config.AdminToken), dep.PlanHandler.Create)
		api.POST("/subscriptions/create", dep.SubscriptionHandler.Create)
		api.GET("/subscriptions/active/:user_id", dep.SubscriptionHandler.GetActive)
		api.POST("/payments/create", dep.PaymentHandler.Create)
		api.POST("/payments/webhook", dep.PaymentHandler.Webhook)
		api.POST("/vpn/issue", dep.VPNHandler.Issue)
	}

	dep.Router.GET("/mock/payments/:external_id", dep.PaymentHandler.ShowMockPaymentPage)
	dep.Router.POST("/mock/payments/:external_id/succeed", dep.PaymentHandler.SimulateSuccess)

	admin := dep.Router.Group("/admin", middleware.AdminAuth(dep.Config.AdminToken))
	{
		admin.GET("", dep.AdminHandler.Dashboard)
		admin.GET("/", dep.AdminHandler.Dashboard)
		admin.GET("/users", dep.AdminHandler.Users)
		admin.GET("/plans", dep.AdminHandler.Plans)
		admin.GET("/subscriptions", dep.AdminHandler.Subscriptions)
		admin.GET("/payments", dep.AdminHandler.Payments)
		admin.GET("/vpn-keys", dep.AdminHandler.VPNKeys)
	}
}
