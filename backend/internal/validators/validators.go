package validators

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"autoservice/backend/internal/dto"
)

func ValidateRegister(req dto.RegisterRequest) error {
	if _, err := mail.ParseAddress(strings.TrimSpace(req.Email)); err != nil {
		return errors.New("invalid email")
	}
	if len(req.Password) < 6 {
		return errors.New("password must contain at least 6 characters")
	}
	if strings.TrimSpace(req.FullName) == "" {
		return errors.New("full_name is required")
	}
	if strings.TrimSpace(req.Phone) == "" {
		return errors.New("phone is required")
	}
	return nil
}

func ValidateLogin(req dto.LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return errors.New("email and password are required")
	}
	return nil
}

func ValidateVehicle(req dto.VehicleCreateRequest) error {
	if strings.TrimSpace(req.Make) == "" || strings.TrimSpace(req.Model) == "" || strings.TrimSpace(req.PlateNumber) == "" {
		return errors.New("make, model and plate_number are required")
	}
	if req.Year < 1900 || req.Year > time.Now().Year()+1 {
		return errors.New("year is out of range")
	}
	return nil
}

func ValidateAppointment(req dto.AppointmentCreateRequest) error {
	if strings.TrimSpace(req.VehicleID) == "" || strings.TrimSpace(req.ServiceID) == "" {
		return errors.New("vehicle_id and service_id are required")
	}
	if strings.TrimSpace(req.StartTime) == "" {
		return errors.New("start_time is required")
	}
	return nil
}

func IsOverlapping(existingStart, existingEnd, newStart, newEnd time.Time) bool {
	return existingStart.Before(newEnd) && existingEnd.After(newStart)
}
