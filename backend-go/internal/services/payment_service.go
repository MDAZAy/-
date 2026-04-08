package services

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		Status:            defaultPaymentStatus(intent.Status),
		Provider:          s.paymentProvider.Name(),
		ExternalPaymentID: intent.ExternalID,
		PaymentURL:        intent.PaymentURL,
		RawResponse:       intent.RawPayload,
		CreatedAt:         time.Now(),
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	if payment.Status == "succeeded" {
		if err := s.fulfillPayment(payment); err != nil {
			return nil, err
		}
	}

	return payment, nil
}

func (s *PaymentService) HandleWebhook(payload []byte, headers http.Header) (*models.Payment, error) {
	event, err := s.paymentProvider.ParseWebhook(payload, headers)
	if err != nil {
		return nil, err
	}

	payment, err := s.repo.FindByExternalID(event.ExternalID)
	if err != nil {
		return nil, err
	}

	if json.Valid(payload) {
		payment.RawResponse = string(payload)
	} else {
		sanitized, _ := json.Marshal(map[string]string{"raw": string(payload)})
		payment.RawResponse = string(sanitized)
	}

	if payment.Status == "succeeded" {
		return payment, s.repo.Save(payment)
	}

	if !event.IsSuccess {
		payment.Status = defaultPaymentStatus(event.Status)
		return payment, s.repo.Save(payment)
	}

	payment.Status = "succeeded"
	if err := s.repo.Save(payment); err != nil {
		return nil, err
	}

	if err := s.fulfillPayment(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) SimulateSuccess(externalID string) (*models.Payment, error) {
	payload, err := json.Marshal(dto.PaymentWebhookRequest{
		Event: "payment.succeeded",
		Object: dto.PaymentWebhookObject{
			ID:     externalID,
			Status: "succeeded",
		},
	})
	if err != nil {
		return nil, err
	}

	return s.HandleWebhook(payload, http.Header{})
}

func (s *PaymentService) ListAll() ([]models.Payment, error) {
	return s.repo.ListAll()
}

func (s *PaymentService) WebhookSuccessResponse() map[string]interface{} {
	return s.paymentProvider.WebhookSuccessResponse()
}

func (s *PaymentService) fulfillPayment(payment *models.Payment) error {
	_, err := s.subscriptionService.Create(dto.CreateSubscriptionRequest{
		UserID: payment.UserID,
		PlanID: payment.PlanID,
	})
	if err != nil {
		return fmt.Errorf("subscription create after payment: %w", err)
	}

	return nil
}

func defaultPaymentStatus(status string) string {
	if status == "" {
		return "pending"
	}
	return status
}
