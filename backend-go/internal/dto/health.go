package dto

type HealthResponse struct {
	Status          string `json:"status"`
	Environment     string `json:"environment"`
	PaymentProvider string `json:"payment_provider"`
	VPNProvider     string `json:"vpn_provider"`
}
