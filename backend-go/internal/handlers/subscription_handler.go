package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"vpn-bot/backend-go/internal/dto"
	"vpn-bot/backend-go/internal/models"
	"vpn-bot/backend-go/internal/services"
)

type SubscriptionHandler struct {
	service *services.SubscriptionService
}

func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var input dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toSubscriptionResponse(subscription))
}

func (h *SubscriptionHandler) GetActive(c *gin.Context) {
	var userID uint
	if _, err := fmt.Sscan(c.Param("user_id"), &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	subscription, err := h.service.GetActive(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "active subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toSubscriptionResponse(subscription))
}

func toSubscriptionResponse(subscription *models.Subscription) dto.SubscriptionResponse {
	return dto.SubscriptionResponse{
		ID:        subscription.ID,
		UserID:    subscription.UserID,
		PlanID:    subscription.PlanID,
		Status:    subscription.Status,
		StartAt:   subscription.StartAt,
		EndAt:     subscription.EndAt,
		CreatedAt: subscription.CreatedAt,
	}
}
