package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"vpn-bot/backend-go/internal/config"
	"vpn-bot/backend-go/internal/dto"
)

var (
	ErrInvalidPaymentWebhook = errors.New("invalid payment webhook")
)

type PaymentIntent struct {
	ExternalID string
	PaymentURL string
	RawPayload string
	Status     string
}

type PaymentWebhookEvent struct {
	ExternalID string
	Status     string
	Event      string
	RawPayload string
	IsSuccess  bool
}

type PaymentProvider interface {
	Name() string
	CreatePayment(userID uint, planID uint, amount float64, returnURL string) (PaymentIntent, error)
	ParseWebhook(payload []byte, headers http.Header) (PaymentWebhookEvent, error)
	WebhookSuccessResponse() map[string]interface{}
}

type MockPaymentProvider struct {
	publicBaseURL string
}

func NewPaymentProvider(cfg config.Config) PaymentProvider {
	switch strings.ToLower(cfg.PaymentProvider) {
	case "cloudpayments":
		return NewCloudPaymentsProvider(cfg)
	default:
		return &MockPaymentProvider{publicBaseURL: cfg.PublicBaseURL}
	}
}

func (p *MockPaymentProvider) Name() string {
	return "mock"
}

func (p *MockPaymentProvider) CreatePayment(userID uint, planID uint, amount float64, _ string) (PaymentIntent, error) {
	externalID := fmt.Sprintf("mockpay_%d_%d_%d", userID, planID, time.Now().UnixNano())
	return PaymentIntent{
		ExternalID: externalID,
		PaymentURL: fmt.Sprintf("%s/mock/payments/%s", p.publicBaseURL, externalID),
		RawPayload: fmt.Sprintf(`{"provider":"mock","amount":%.2f}`, amount),
		Status:     "pending",
	}, nil
}

func (p *MockPaymentProvider) ParseWebhook(payload []byte, _ http.Header) (PaymentWebhookEvent, error) {
	var input dto.PaymentWebhookRequest
	if err := json.Unmarshal(payload, &input); err != nil {
		return PaymentWebhookEvent{}, fmt.Errorf("%w: decode mock webhook: %v", ErrInvalidPaymentWebhook, err)
	}

	return PaymentWebhookEvent{
		ExternalID: input.Object.ID,
		Status:     input.Object.Status,
		Event:      input.Event,
		RawPayload: string(payload),
		IsSuccess:  input.Object.Status == "succeeded" || input.Event == "payment.succeeded",
	}, nil
}

func (p *MockPaymentProvider) WebhookSuccessResponse() map[string]interface{} {
	return map[string]interface{}{
		"code":    0,
		"message": "ok",
	}
}
