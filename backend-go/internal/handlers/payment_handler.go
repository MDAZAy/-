package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"vpn-bot/backend-go/internal/dto"
	"vpn-bot/backend-go/internal/models"
	"vpn-bot/backend-go/internal/services"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var input dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toPaymentResponse(payment))
}

func (h *PaymentHandler) Webhook(c *gin.Context) {
	var input dto.PaymentWebhookRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.HandleWebhook(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toPaymentResponse(payment))
}

func (h *PaymentHandler) ShowMockPaymentPage(c *gin.Context) {
	c.HTML(http.StatusOK, "mock_payment.tmpl", gin.H{
		"Title":      "Mock Payment",
		"ExternalID": c.Param("external_id"),
	})
}

func (h *PaymentHandler) SimulateSuccess(c *gin.Context) {
	payment, err := h.service.SimulateSuccess(c.Param("external_id"))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "mock_payment_result.tmpl", gin.H{
			"Title":   "Payment Error",
			"Success": false,
			"Message": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "mock_payment_result.tmpl", gin.H{
		"Title":   "Payment Complete",
		"Success": true,
		"Message": "Платёж отмечен как успешный, подписка создана.",
		"Payment": toPaymentResponse(payment),
	})
}

func toPaymentResponse(payment *models.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:                payment.ID,
		UserID:            payment.UserID,
		PlanID:            payment.PlanID,
		Amount:            payment.Amount,
		Currency:          payment.Currency,
		Status:            payment.Status,
		Provider:          payment.Provider,
		ExternalPaymentID: payment.ExternalPaymentID,
		PaymentURL:        payment.PaymentURL,
		CreatedAt:         payment.CreatedAt,
	}
}
