package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"vpn-bot/backend-go/internal/models"
	"vpn-bot/backend-go/internal/providers"
	"vpn-bot/backend-go/internal/repositories"
)

type VPNService struct {
	subscriptionRepo *repositories.SubscriptionRepository
	keyRepo          *repositories.VPNKeyRepository
	provider         providers.VPNProvider
}

func NewVPNService(
	subscriptionRepo *repositories.SubscriptionRepository,
	keyRepo *repositories.VPNKeyRepository,
	provider providers.VPNProvider,
) *VPNService {
	return &VPNService{
		subscriptionRepo: subscriptionRepo,
		keyRepo:          keyRepo,
		provider:         provider,
	}
}

func (s *VPNService) IssueKey(userID uint) (*models.VPNKey, error) {
	subscription, err := s.subscriptionRepo.GetActiveByUser(userID, time.Now())
	if err != nil {
		return nil, err
	}

	key, err := s.keyRepo.FindActiveByUser(userID)
	if err == nil {
		return key, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	issued, err := s.provider.IssueKey(userID, subscription.EndAt)
	if err != nil {
		return nil, err
	}

	record := &models.VPNKey{
		UserID:           userID,
		Provider:         s.provider.Name(),
		ExternalClientID: issued.ExternalClientID,
		KeyName:          issued.KeyName,
		AccessURL:        issued.AccessURL,
		ConfigJSON:       issued.ConfigJSON,
		IsActive:         true,
		ExpiresAt:        &subscription.EndAt,
		CreatedAt:        time.Now(),
	}

	return record, s.keyRepo.Create(record)
}

func (s *VPNService) ListAll() ([]models.VPNKey, error) {
	return s.keyRepo.ListAll()
}
