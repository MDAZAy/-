package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vpn-bot/backend-go/internal/services"
)

type AdminHandler struct {
	userService         *services.UserService
	planService         *services.PlanService
	subscriptionService *services.SubscriptionService
	paymentService      *services.PaymentService
	vpnService          *services.VPNService
}

func NewAdminHandler(
	userService *services.UserService,
	planService *services.PlanService,
	subscriptionService *services.SubscriptionService,
	paymentService *services.PaymentService,
	vpnService *services.VPNService,
) *AdminHandler {
	return &AdminHandler{
		userService:         userService,
		planService:         planService,
		subscriptionService: subscriptionService,
		paymentService:      paymentService,
		vpnService:          vpnService,
	}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	users, _ := h.userService.ListAll()
	plans, _ := h.planService.ListAll()
	subscriptions, _ := h.subscriptionService.ListAll()
	payments, _ := h.paymentService.ListAll()
	keys, _ := h.vpnService.ListAll()

	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"Title":         "Dashboard",
		"UsersCount":    len(users),
		"PlansCount":    len(plans),
		"SubsCount":     len(subscriptions),
		"PaymentsCount": len(payments),
		"KeysCount":     len(keys),
	})
}

func (h *AdminHandler) Users(c *gin.Context) {
	users, _ := h.userService.ListAll()
	c.HTML(http.StatusOK, "users.tmpl", gin.H{"Title": "Users", "Items": users})
}

func (h *AdminHandler) Plans(c *gin.Context) {
	plans, _ := h.planService.ListAll()
	c.HTML(http.StatusOK, "plans.tmpl", gin.H{"Title": "Plans", "Items": plans})
}

func (h *AdminHandler) Subscriptions(c *gin.Context) {
	subscriptions, _ := h.subscriptionService.ListAll()
	c.HTML(http.StatusOK, "subscriptions.tmpl", gin.H{"Title": "Subscriptions", "Items": subscriptions})
}

func (h *AdminHandler) Payments(c *gin.Context) {
	payments, _ := h.paymentService.ListAll()
	c.HTML(http.StatusOK, "payments.tmpl", gin.H{"Title": "Payments", "Items": payments})
}

func (h *AdminHandler) VPNKeys(c *gin.Context) {
	keys, _ := h.vpnService.ListAll()
	c.HTML(http.StatusOK, "vpn_keys.tmpl", gin.H{"Title": "VPN Keys", "Items": keys})
}
