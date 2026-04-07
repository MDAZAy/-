package repositories

import (
	"gorm.io/gorm"

	"vpn-bot/backend-go/internal/models"
)

type VPNKeyRepository struct {
	db *gorm.DB
}

func NewVPNKeyRepository(db *gorm.DB) *VPNKeyRepository {
	return &VPNKeyRepository{db: db}
}

func (r *VPNKeyRepository) FindActiveByUser(userID uint) (*models.VPNKey, error) {
	var key models.VPNKey
	if err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at desc").
		First(&key).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *VPNKeyRepository) Create(key *models.VPNKey) error {
	return r.db.Create(key).Error
}

func (r *VPNKeyRepository) Save(key *models.VPNKey) error {
	return r.db.Save(key).Error
}

func (r *VPNKeyRepository) DeactivateAllByUser(userID uint) error {
	return r.db.Model(&models.VPNKey{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false).Error
}

func (r *VPNKeyRepository) ListAll() ([]models.VPNKey, error) {
	var keys []models.VPNKey
	err := r.db.Order("created_at desc").Find(&keys).Error
	return keys, err
}
