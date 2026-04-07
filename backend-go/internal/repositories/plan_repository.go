package repositories

import (
	"gorm.io/gorm"

	"vpn-bot/backend-go/internal/models"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) ListActive() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Where("is_active = ?", true).Order("price asc").Find(&plans).Error
	return plans, err
}

func (r *PlanRepository) ListAll() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Order("created_at desc").Find(&plans).Error
	return plans, err
}

func (r *PlanRepository) GetByID(id uint) (*models.Plan, error) {
	var plan models.Plan
	if err := r.db.First(&plan, id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) Create(plan *models.Plan) error {
	return r.db.Create(plan).Error
}
