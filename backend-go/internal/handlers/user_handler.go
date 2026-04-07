package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vpn-bot/backend-go/internal/dto"
	"vpn-bot/backend-go/internal/models"
	"vpn-bot/backend-go/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Ensure(c *gin.Context) {
	var input dto.EnsureUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.EnsureUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         user.ID,
		TelegramID: user.TelegramID,
		Username:   user.Username,
		FullName:   user.FullName,
		IsAdmin:    user.IsAdmin,
		IsBlocked:  user.IsBlocked,
		CreatedAt:  user.CreatedAt,
	}
}
