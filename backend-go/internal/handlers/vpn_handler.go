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

type VPNHandler struct {
	service *services.VPNService
}

func NewVPNHandler(service *services.VPNService) *VPNHandler {
	return &VPNHandler{service: service}
}

func (h *VPNHandler) Issue(c *gin.Context) {
	var input dto.IssueVPNKeyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := h.service.IssueKey(input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "active subscription required"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toVPNResponse(key))
}

func toVPNResponse(key *models.VPNKey) dto.VPNKeyResponse {
	return dto.VPNKeyResponse{
		ID:               key.ID,
		UserID:           key.UserID,
		Provider:         key.Provider,
		ExternalClientID: key.ExternalClientID,
		KeyName:          key.KeyName,
		AccessURL:        key.AccessURL,
		ConfigJSON:       key.ConfigJSON,
		IsActive:         key.IsActive,
		ExpiresAt:        key.ExpiresAt,
		CreatedAt:        key.CreatedAt,
	}
}
