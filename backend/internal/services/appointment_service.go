package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"autoservice/backend/internal/cache"
	"autoservice/backend/internal/config"
	"autoservice/backend/internal/dto"
	"autoservice/backend/internal/models"
	"autoservice/backend/internal/repositories"
	"autoservice/backend/internal/validators"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentService struct {
	repo  *repositories.Repository
	cache *cache.TTLCache
	tz    *time.Location
}

func NewAppointmentService(repo *repositories.Repository, cfg config.Config) *AppointmentService {
	location, err := time.LoadLocation(cfg.DisplayTimezone)
	if err != nil {
		location = time.UTC
	}

	return &AppointmentService{
		repo:  repo,
		cache: cache.NewTTLCache(),
		tz:    location,
	}
}

func (s *AppointmentService) CreateAppointment(userID string, req dto.AppointmentCreateRequest, idempotencyKey string, meta RequestMeta) (*dto.AppointmentResponse, *AppError) {
	if err := validators.ValidateAppointment(req); err != nil {
		return nil, NewError(http.StatusBadRequest, "validation_failed", err.Error())
	}
	if strings.TrimSpace(idempotencyKey) == "" {
		return nil, NewError(http.StatusBadRequest, "idempotency_key_required", "Idempotency-Key header is required")
	}

	startUTC, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, NewError(http.StatusBadRequest, "invalid_start_time", "start_time must be RFC3339 in UTC")
	}
	startUTC = startUTC.UTC()
	if startUTC.Before(time.Now().UTC()) {
		return nil, NewError(http.StatusBadRequest, "appointment_in_past", "appointment cannot be created in the past")
	}

	tx := s.repo.DB().Begin()
	if tx.Error != nil {
		return nil, NewError(http.StatusInternalServerError, "transaction_start_failed", "failed to start transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if existing, err := s.repo.FindAppointmentByIdempotencyTx(tx, userID, idempotencyKey); err == nil {
		if commitErr := tx.Commit().Error; commitErr != nil {
			return nil, NewError(http.StatusInternalServerError, "transaction_commit_failed", "failed to finish idempotent request")
		}
		response := s.mapAppointment(*existing)
		return &response, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, NewError(http.StatusInternalServerError, "idempotency_check_failed", "failed to validate idempotency key")
	}

	vehicle, err := s.repo.FindVehicleForUserTx(tx, userID, req.VehicleID)
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, "vehicle_not_found", "vehicle not found for user")
	}

	service, err := s.repo.FindServiceTx(tx, req.ServiceID)
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, "service_not_found", "service not found")
	}

	endUTC := startUTC.Add(time.Duration(service.DurationMinutes) * time.Minute)
	if req.EndTime != nil && strings.TrimSpace(*req.EndTime) != "" {
		parsedEnd, parseErr := time.Parse(time.RFC3339, *req.EndTime)
		if parseErr != nil {
			tx.Rollback()
			return nil, NewError(http.StatusBadRequest, "invalid_end_time", "end_time must be RFC3339 in UTC")
		}
		endUTC = parsedEnd.UTC()
	}
	if !endUTC.After(startUTC) {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, "invalid_time_range", "end_time must be later than start_time")
	}

	startLocal := startUTC.In(s.tz)
	endLocal := endUTC.In(s.tz)
	if startLocal.Format("2006-01-02") != endLocal.Format("2006-01-02") {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, "cross_day_appointment", "appointment must stay within one working day")
	}

	workingHour, err := s.repo.FindWorkingHourTx(tx, int(startLocal.Weekday()))
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, "working_hours_not_found", "working hours are not configured")
	}
	if !workingHour.IsWorking || !withinWorkingHours(startLocal, endLocal, workingHour.StartTime, workingHour.EndTime) {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, "outside_working_hours", "selected time is outside working hours")
	}

	isHoliday, err := s.repo.IsHolidayTx(tx, startLocal)
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusInternalServerError, "holiday_check_failed", "failed to validate holidays")
	}
	if isHoliday {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, "holiday_blocked", "service is closed on holidays")
	}

	mechanic, err := s.repo.FindFreeMechanicTx(tx, startUTC, endUTC)
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusConflict, "slot_unavailable", "no available mechanic for selected time")
	}

	status, err := s.repo.FindStatusByCodeTx(tx, "scheduled")
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusInternalServerError, "status_not_found", "appointment status is not configured")
	}

	appointment := models.Appointment{
		UserID:             userID,
		VehicleID:          vehicle.ID,
		ServiceID:          service.ID,
		StatusID:           status.ID,
		MechanicID:         mechanic.ID,
		StartTime:          startUTC,
		EndTime:            endUTC,
		ConfirmationNumber: buildConfirmationNumber(),
		IdempotencyKey:     strings.TrimSpace(idempotencyKey),
		Notes:              strings.TrimSpace(req.Notes),
	}

	if err := s.repo.CreateAppointmentTx(tx, &appointment); err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusInternalServerError, "appointment_create_failed", "failed to create appointment")
	}

	if err := s.repo.CreateAppointmentHistoryTx(tx, &models.AppointmentHistory{
		AppointmentID:   appointment.ID,
		StatusID:        status.ID,
		ChangedByUserID: &userID,
		Comment:         "Appointment created",
	}); err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusInternalServerError, "appointment_history_failed", "failed to save appointment history")
	}

	if err := s.repo.CreateNotificationTx(tx, &models.Notification{
		UserID:        userID,
		AppointmentID: &appointment.ID,
		Type:          "appointment_created",
		Message:       fmt.Sprintf("Your appointment %s has been created", appointment.ConfirmationNumber),
	}); err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusInternalServerError, "notification_failed", "failed to save notification")
	}

	if err := s.repo.CreateAuditLogTx(tx, &models.AuditLog{
		UserID:      &userID,
		Action:      "appointment_create",
		Entity:      "appointment",
		EntityID:    appointment.ID,
		IPAddress:   meta.IPAddress,
		Metadata:    fmt.Sprintf(`{"vehicle_id":"%s","service_id":"%s","mechanic_id":"%s"}`, vehicle.ID, service.ID, mechanic.ID),
		Description: "Appointment created via online booking",
	}); err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusInternalServerError, "audit_log_failed", "failed to save audit log")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, NewError(http.StatusInternalServerError, "transaction_commit_failed", "failed to save appointment")
	}

	appointment.Vehicle = *vehicle
	appointment.Service = *service
	appointment.Status = *status
	appointment.Mechanic = *mechanic
	s.cache.DeletePrefix("slots:")

	response := s.mapAppointment(appointment)
	return &response, nil
}

func (s *AppointmentService) ListMyAppointments(userID string) ([]dto.AppointmentResponse, *AppError) {
	appointments, err := s.repo.ListUserAppointments(userID)
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "appointments_load_failed", "failed to load appointments")
	}
	return s.mapAppointments(appointments), nil
}

func (s *AppointmentService) ListAllAppointments() ([]dto.AppointmentResponse, *AppError) {
	appointments, err := s.repo.ListAllAppointments()
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "appointments_load_failed", "failed to load appointments")
	}
	return s.mapAppointments(appointments), nil
}

func (s *AppointmentService) AvailableSlots(dateRaw, serviceID string) (*dto.AvailableSlotsResponse, *AppError) {
	cacheKey := fmt.Sprintf("slots:%s:%s", dateRaw, serviceID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		result := cached.(dto.AvailableSlotsResponse)
		return &result, nil
	}

	dayLocal, err := time.ParseInLocation("2006-01-02", dateRaw, s.tz)
	if err != nil {
		return nil, NewError(http.StatusBadRequest, "invalid_date", "date must be YYYY-MM-DD")
	}

	service, err := s.repo.GetService(serviceID)
	if err != nil {
		return nil, NewError(http.StatusBadRequest, "service_not_found", "service not found")
	}

	tx := s.repo.DB().Begin()
	if tx.Error != nil {
		return nil, NewError(http.StatusInternalServerError, "transaction_start_failed", "failed to validate schedule")
	}

	workingHour, err := s.repo.FindWorkingHourTx(tx, int(dayLocal.Weekday()))
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusBadRequest, "working_hours_not_found", "working hours are not configured")
	}
	isHoliday, err := s.repo.IsHolidayTx(tx, dayLocal)
	if err != nil {
		tx.Rollback()
		return nil, NewError(http.StatusInternalServerError, "holiday_check_failed", "failed to validate holidays")
	}
	if err := tx.Commit().Error; err != nil {
		return nil, NewError(http.StatusInternalServerError, "transaction_commit_failed", "failed to read schedule")
	}

	response := dto.AvailableSlotsResponse{
		Date:     dateRaw,
		Timezone: s.tz.String(),
		Slots:    []dto.SlotResponse{},
	}
	if !workingHour.IsWorking || isHoliday {
		s.cache.Set(cacheKey, response, 5*time.Minute)
		return &response, nil
	}

	mechanicsCount, err := s.repo.CountActiveMechanics()
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "mechanics_count_failed", "failed to count mechanics")
	}
	if mechanicsCount == 0 {
		return &response, nil
	}

	openLocal := mustClock(dayLocal, workingHour.StartTime, s.tz)
	closeLocal := mustClock(dayLocal, workingHour.EndTime, s.tz)
	dayAppointments, err := s.repo.ListDayAppointments(openLocal.UTC(), closeLocal.UTC())
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "appointments_load_failed", "failed to calculate available slots")
	}

	step := 30 * time.Minute
	duration := time.Duration(service.DurationMinutes) * time.Minute
	for slotStart := openLocal; !slotStart.Add(duration).After(closeLocal); slotStart = slotStart.Add(step) {
		slotEnd := slotStart.Add(duration)
		busy := 0
		for _, item := range dayAppointments {
			if validators.IsOverlapping(item.StartTime, item.EndTime, slotStart.UTC(), slotEnd.UTC()) {
				busy++
			}
		}
		if int64(busy) >= mechanicsCount {
			continue
		}
		response.Slots = append(response.Slots, dto.SlotResponse{
			StartTimeUTC:   slotStart.UTC().Format(time.RFC3339),
			EndTimeUTC:     slotEnd.UTC().Format(time.RFC3339),
			StartTimeLocal: slotStart.Format(time.RFC3339),
			EndTimeLocal:   slotEnd.Format(time.RFC3339),
		})
	}

	s.cache.Set(cacheKey, response, 5*time.Minute)
	return &response, nil
}

func (s *AppointmentService) mapAppointments(items []models.Appointment) []dto.AppointmentResponse {
	result := make([]dto.AppointmentResponse, 0, len(items))
	for _, item := range items {
		result = append(result, s.mapAppointment(item))
	}
	return result
}

func (s *AppointmentService) mapAppointment(item models.Appointment) dto.AppointmentResponse {
	startLocal := item.StartTime.In(s.tz)
	endLocal := item.EndTime.In(s.tz)
	return dto.AppointmentResponse{
		ID:                 item.ID,
		ConfirmationNumber: item.ConfirmationNumber,
		Status:             item.Status.Code,
		Service: dto.ServiceResponse{
			ID:              item.Service.ID,
			CategoryID:      item.Service.CategoryID,
			CategoryName:    item.Service.Category.Name,
			Name:            item.Service.Name,
			Description:     item.Service.Description,
			DurationMinutes: item.Service.DurationMinutes,
			Price:           item.Service.Price,
		},
		Vehicle: dto.VehicleResponse{
			ID:          item.Vehicle.ID,
			Make:        item.Vehicle.Make,
			Model:       item.Vehicle.Model,
			Year:        item.Vehicle.Year,
			PlateNumber: item.Vehicle.PlateNumber,
			Color:       item.Vehicle.Color,
			VIN:         item.Vehicle.VIN,
		},
		MechanicName:   item.Mechanic.FullName,
		StartTimeUTC:   item.StartTime.UTC().Format(time.RFC3339),
		EndTimeUTC:     item.EndTime.UTC().Format(time.RFC3339),
		StartTimeLocal: startLocal.Format(time.RFC3339),
		EndTimeLocal:   endLocal.Format(time.RFC3339),
		Notes:          item.Notes,
	}
}

func withinWorkingHours(startLocal, endLocal time.Time, openClock, closeClock string) bool {
	openTime := mustClock(startLocal, openClock, startLocal.Location())
	closeTime := mustClock(startLocal, closeClock, startLocal.Location())
	return !startLocal.Before(openTime) && !endLocal.After(closeTime)
}

func mustClock(day time.Time, clock string, loc *time.Location) time.Time {
	parsed, _ := time.ParseInLocation("2006-01-02 15:04", day.Format("2006-01-02")+" "+clock, loc)
	return parsed
}

func buildConfirmationNumber() string {
	return "AS-" + strings.ToUpper(uuid.NewString()[:8])
}
