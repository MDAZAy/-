package dto

import "time"

type IssueVPNKeyRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

type VPNKeyResponse struct {
	ID               uint       `json:"id"`
	UserID           uint       `json:"user_id"`
	Provider         string     `json:"provider"`
	ExternalClientID string     `json:"external_client_id"`
	KeyName          string     `json:"key_name"`
	AccessURL        string     `json:"access_url"`
	ConfigJSON       string     `json:"config_json"`
	IsActive         bool       `json:"is_active"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}
