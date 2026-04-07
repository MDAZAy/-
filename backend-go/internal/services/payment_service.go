package services

import (
	"encoding/json"
	"fmt"
	"time"

	"vpn-bot/backend-go/internal/dto"
	"vpn-bot/backend-go/internal/models"
	"vpn-bot/backend-go/internal/providers"
	"vpn-bot/backend-go/internal/repositories"
)

type PaymentService struct {
	repo                *repositories.PaymentRepository
	planRepo            *repositories.PlanRepository
	subscriptionService *SubscriptionService
	paymentProvider     providers.PaymentProvider
}

func NewPaymentService(
	repo *repositories.PaymentRepository,
	planRepo *repositories.PlanRepository,
	subscriptionService *SubscriptionService,
	paymentProvider providers.PaymentProvider,
) *PaymentService {
	return &PaymentService{
		repo:                repo,
		planRepo:            planRepo,
		subscriptionService: subscriptionService,
		paymentProvider:     paymentProvider,
	}
}

func (s *PaymentService) Create(input dto.CreatePaymentRequest) (*models.Payment, error) {
	plan, err := s.planRepo.GetByID(input.PlanID)
	if err != nil {
		return nil, err
	}

	intent, err := s.paymentProvider.CreatePayment(input.UserID, input.PlanID, plan.Price, input.ReturnURL)
	if err != nil {
		return nil, err
	}

	payment := &models.Payment{
		UserID:            input.UserID,
		PlanID:            input.PlanID,
		Amount:            plan.Price,
		Currency:          "RUB",
		Status:            "pending",
		Provider:          s.paymentProvider.Name(),
		ExternalPaymentID: intent.ExternalID,
		PaymentURL:        intent.PaymentURL,
		RawResponse:       intent.RawPayload,
		CreatedAt:         time.Now(),
	}

	return payment, s.repo.Create(payment)
}

func (s *PaymentService) HandleWebhook(input dto.PaymentWebhookRequest) (*models.Payment, error) {
	payment, err := s.repo.FindByExternalID(input.Object.ID)
	if err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"event":  input.Event,
		"object": input.Object,
	})
	payment.RawResponse = string(payload)

	if payment.Status == "succeeded" {
		return payment, s.repo.Save(payment)
	}

	if input.Object.Status != "succeeded" && input.Event != "payment.succeeded" {
		payment.Status = input.Object.Status
		return payment, s.repo.Save(payment)
	}

	payment.Status = "succeeded"
	if err := s.repo.Save(payment); err != nil {
		return nil, err
	}

	_, err = s.subscriptionService.Create(dto.CreateSubscriptionRequest{
		UserID: payment.UserID,
		PlanID: payment.PlanID,
	})
	if err != nil {
		return nil, fmt.Errorf("subscription create after payment: %w", err)
	}

	return payment, nil
}

func (s *PaymentService) SimulateSuccess(externalID string) (*models.Payment, error) {
	return s.HandleWebhook(dto.PaymentWebhookRequest{
		Event: "payment.succeeded",
		Object: dto.PaymentWebhookObject{
			ID:     externalID,
			Status: "succeeded",
		},
	})
}

func (s *PaymentService) ListAll() ([]models.Payment, error) {
	return s.repo.ListAll()
}
