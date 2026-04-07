package providers

import (
	"fmt"
	"time"

	"vpn-bot/backend-go/internal/config"
)

type PaymentIntent struct {
	ExternalID string
	PaymentURL string
	RawPayload string
}

type PaymentProvider interface {
	Name() string
	CreatePayment(userID uint, planID uint, amount float64, returnURL string) (PaymentIntent, error)
}

type MockPaymentProvider struct {
	publicBaseURL string
}

func NewPaymentProvider(cfg config.Config) PaymentProvider {
	return &MockPaymentProvider{publicBaseURL: cfg.PublicBaseURL}
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
	}, nil
}
