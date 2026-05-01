package handlers

import (
	"autoservice/backend/internal/dto"
	"autoservice/backend/internal/middleware"
	"autoservice/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) listCategories(c *gin.Context) {
	data, appErr := h.catalogService.ListCategories()
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func (h *HTTPHandler) listServices(c *gin.Context) {
	data, appErr := h.catalogService.ListServices()
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func (h *HTTPHandler) createVehicle(c *gin.Context) {
	var req dto.VehicleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		failure(c, services.NewError(400, "invalid_payload", "invalid request body"))
		return
	}

	userID := c.GetString(middleware.ContextUserID)
	data, appErr := h.catalogService.CreateVehicle(userID, req)
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func (h *HTTPHandler) listVehicles(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserID)
	data, appErr := h.catalogService.ListUserVehicles(userID)
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func (h *HTTPHandler) me(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserID)
	data, appErr := h.catalogService.Profile(userID)
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func (h *HTTPHandler) dashboard(c *gin.Context) {
	data, appErr := h.catalogService.AdminDashboard()
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}
