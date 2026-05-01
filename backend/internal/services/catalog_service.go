package services

import (
	"fmt"
	"net/http"
	"time"

	"autoservice/backend/internal/cache"
	"autoservice/backend/internal/config"
	"autoservice/backend/internal/dto"
	"autoservice/backend/internal/models"
	"autoservice/backend/internal/repositories"
	"autoservice/backend/internal/validators"
)

type CatalogService struct {
	repo  *repositories.Repository
	cache *cache.TTLCache
}

func NewCatalogService(repo *repositories.Repository, _ config.Config) *CatalogService {
	return &CatalogService{
		repo:  repo,
		cache: cache.NewTTLCache(),
	}
}

func (s *CatalogService) ListCategories() ([]dto.CategoryResponse, *AppError) {
	if cached, ok := s.cache.Get("categories"); ok {
		return cached.([]dto.CategoryResponse), nil
	}

	categories, err := s.repo.ListActiveCategories()
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "categories_load_failed", "failed to load categories")
	}

	result := make([]dto.CategoryResponse, 0, len(categories))
	for _, item := range categories {
		result = append(result, dto.CategoryResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		})
	}

	s.cache.Set("categories", result, 15*time.Minute)
	return result, nil
}

func (s *CatalogService) ListServices() ([]dto.ServiceResponse, *AppError) {
	if cached, ok := s.cache.Get("services"); ok {
		return cached.([]dto.ServiceResponse), nil
	}

	servicesList, err := s.repo.ListActiveServices()
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "services_load_failed", "failed to load services")
	}

	result := make([]dto.ServiceResponse, 0, len(servicesList))
	for _, item := range servicesList {
		result = append(result, dto.ServiceResponse{
			ID:              item.ID,
			CategoryID:      item.CategoryID,
			CategoryName:    item.Category.Name,
			Name:            item.Name,
			Description:     item.Description,
			DurationMinutes: item.DurationMinutes,
			Price:           item.Price,
		})
	}

	s.cache.Set("services", result, 15*time.Minute)
	return result, nil
}

func (s *CatalogService) CreateVehicle(userID string, req dto.VehicleCreateRequest) (*dto.VehicleResponse, *AppError) {
	if err := validators.ValidateVehicle(req); err != nil {
		return nil, NewError(http.StatusBadRequest, "validation_failed", err.Error())
	}

	vehicle := models.Vehicle{
		UserID:      userID,
		Make:        req.Make,
		Model:       req.Model,
		Year:        req.Year,
		PlateNumber: req.PlateNumber,
		Color:       req.Color,
		VIN:         req.VIN,
	}

	if err := s.repo.CreateVehicle(&vehicle); err != nil {
		return nil, NewError(http.StatusConflict, "vehicle_create_failed", "failed to create vehicle")
	}

	return &dto.VehicleResponse{
		ID:          vehicle.ID,
		Make:        vehicle.Make,
		Model:       vehicle.Model,
		Year:        vehicle.Year,
		PlateNumber: vehicle.PlateNumber,
		Color:       vehicle.Color,
		VIN:         vehicle.VIN,
	}, nil
}

func (s *CatalogService) ListUserVehicles(userID string) ([]dto.VehicleResponse, *AppError) {
	vehicles, err := s.repo.ListUserVehicles(userID)
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "vehicles_load_failed", "failed to load vehicles")
	}

	result := make([]dto.VehicleResponse, 0, len(vehicles))
	for _, item := range vehicles {
		result = append(result, dto.VehicleResponse{
			ID:          item.ID,
			Make:        item.Make,
			Model:       item.Model,
			Year:        item.Year,
			PlateNumber: item.PlateNumber,
			Color:       item.Color,
			VIN:         item.VIN,
		})
	}
	return result, nil
}

func (s *CatalogService) Profile(userID string) (*dto.ProfileResponse, *AppError) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, NewError(http.StatusNotFound, "user_not_found", "user not found")
	}

	return &dto.ProfileResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Role:     user.Role.Name,
	}, nil
}

func (s *CatalogService) AdminDashboard() (*dto.DashboardResponse, *AppError) {
	counts, err := s.repo.GetDashboardCounts()
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "dashboard_load_failed", "failed to load dashboard")
	}

	return &dto.DashboardResponse{
		UsersCount:        counts["users"],
		VehiclesCount:     counts["vehicles"],
		AppointmentsCount: counts["appointments"],
		MechanicsCount:    counts["mechanics"],
	}, nil
}

func (s *CatalogService) InvalidateScheduleCache(date string, serviceID string) {
	s.cache.Delete(fmt.Sprintf("slots:%s:%s", date, serviceID))
}
