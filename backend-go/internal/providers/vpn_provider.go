package providers

import (
	"fmt"
	"time"

	"vpn-bot/backend-go/internal/config"
)

type IssuedVPNKey struct {
	ExternalClientID string
	KeyName          string
	AccessURL        string
	ConfigJSON       string
}

type VPNProvider interface {
	Name() string
	IssueKey(userID uint, endAt time.Time) (IssuedVPNKey, error)
	DeactivateKey(externalClientID string) error
}

type MockVPNProvider struct {
	publicBaseURL string
}

func NewVPNProvider(cfg config.Config) VPNProvider {
	return &MockVPNProvider{publicBaseURL: cfg.PublicBaseURL}
}

func (p *MockVPNProvider) Name() string {
	return "mock"
}

func (p *MockVPNProvider) IssueKey(userID uint, endAt time.Time) (IssuedVPNKey, error) {
	externalID := fmt.Sprintf("mockvpn_%d_%d", userID, time.Now().UnixNano())
	return IssuedVPNKey{
		ExternalClientID: externalID,
		KeyName:          fmt.Sprintf("vpn-user-%d", userID),
		AccessURL:        fmt.Sprintf("vpn://mock/access/%s", externalID),
		ConfigJSON:       fmt.Sprintf(`{"provider":"mock","expires_at":"%s"}`, endAt.Format(time.RFC3339)),
	}, nil
}

func (p *MockVPNProvider) DeactivateKey(string) error {
	return nil
}
