package handlers

import (
	"autoservice/backend/internal/dto"
	"autoservice/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		failure(c, services.NewError(400, "invalid_payload", "invalid request body"))
		return
	}

	response, appErr := h.authService.Register(req, requestMeta(c))
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, response)
}

func (h *HTTPHandler) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		failure(c, services.NewError(400, "invalid_payload", "invalid request body"))
		return
	}

	response, appErr := h.authService.Login(req, requestMeta(c))
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, response)
}

func (h *HTTPHandler) refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		failure(c, services.NewError(400, "invalid_payload", "invalid request body"))
		return
	}

	response, appErr := h.authService.Refresh(req, requestMeta(c))
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, response)
}

func (h *HTTPHandler) logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		failure(c, services.NewError(400, "invalid_payload", "invalid request body"))
		return
	}
	if appErr := h.authService.Logout(req); appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, gin.H{"logged_out": true})
}
