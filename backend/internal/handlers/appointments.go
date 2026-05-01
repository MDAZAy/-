package handlers

import (
	"autoservice/backend/internal/dto"
	"autoservice/backend/internal/middleware"
	"autoservice/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) createAppointment(c *gin.Context) {
	var req dto.AppointmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		failure(c, services.NewError(400, "invalid_payload", "invalid request body"))
		return
	}

	userID := c.GetString(middleware.ContextUserID)
	data, appErr := h.appointmentService.CreateAppointment(userID, req, c.GetHeader("Idempotency-Key"), requestMeta(c))
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func (h *HTTPHandler) listMyAppointments(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserID)
	data, appErr := h.appointmentService.ListMyAppointments(userID)
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func (h *HTTPHandler) listAllAppointments(c *gin.Context) {
	data, appErr := h.appointmentService.ListAllAppointments()
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func (h *HTTPHandler) availableSlots(c *gin.Context) {
	date := c.Query("date")
	serviceID := c.Query("service_id")
	data, appErr := h.appointmentService.AvailableSlots(date, serviceID)
	if appErr != nil {
		failure(c, appErr)
		return
	}
	success(c, data)
}

func requestMeta(c *gin.Context) services.RequestMeta {
	return services.RequestMeta{
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
}
