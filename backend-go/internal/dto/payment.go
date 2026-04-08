package dto

import "time"

type CreatePaymentRequest struct {
	UserID    uint   `json:"user_id" binding:"required"`
	PlanID    uint   `json:"plan_id" binding:"required"`
	ReturnURL string `json:"return_url"`
}

type PaymentWebhookRequest struct {
	Event  string               `json:"event"`
	Object PaymentWebhookObject `json:"object"`
}

type PaymentWebhookObject struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type PaymentResponse struct {
	ID                uint      `json:"id"`
	UserID            uint      `json:"user_id"`
	PlanID            uint      `json:"plan_id"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	Status            string    `json:"status"`
	Provider          string    `json:"provider"`
	ExternalPaymentID string    `json:"external_payment_id"`
	PaymentURL        string    `json:"payment_url"`
	CreatedAt         time.Time `json:"created_at"`
}

type PaymentReturnPageData struct {
	Title       string
	Success     bool
	Message     string
	SupportURL  string
	TelegramBot string
}
