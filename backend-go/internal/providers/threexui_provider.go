package providers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"vpn-bot/backend-go/internal/config"
)

type ThreeXUIVPNProvider struct {
	panelURL    string
	username    string
	password    string
	inboundID   int
	publicHost  string
	publicPort  string
	serverName  string
	publicKey   string
	shortID     string
	flow        string
	fingerprint string
	client      *http.Client
}

type threeXUILoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type threeXUIBaseResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

type threeXUIClientSettings struct {
	Clients []threeXUIClient `json:"clients"`
}

type threeXUIClient struct {
	ID         string `json:"id"`
	Flow       string `json:"flow,omitempty"`
	Email      string `json:"email"`
	LimitIP    int    `json:"limitIp"`
	TotalGB    int64  `json:"totalGB"`
	ExpiryTime int64  `json:"expiryTime"`
	Enable     bool   `json:"enable"`
	TGID       string `json:"tgId"`
	SubID      string `json:"subId,omitempty"`
	Comment    string `json:"comment,omitempty"`
	Reset      int    `json:"reset"`
}

type threeXUIAddClientRequest struct {
	ID       int    `json:"id"`
	Settings string `json:"settings"`
}

func NewThreeXUIVPNProvider(cfg config.Config) VPNProvider {
	jar, _ := cookiejar.New(nil)

	return &ThreeXUIVPNProvider{
		panelURL:    strings.TrimRight(cfg.VPNProviderURL, "/"),
		username:    cfg.VPNProviderUsername,
		password:    cfg.VPNProviderPassword,
		inboundID:   cfg.VPNProviderInboundID,
		publicHost:  cfg.VPNPublicHost,
		publicPort:  cfg.VPNPublicPort,
		serverName:  cfg.VPNRealityServerName,
		publicKey:   cfg.VPNRealityPublicKey,
		shortID:     cfg.VPNRealityShortID,
		flow:        cfg.VPNFlow,
		fingerprint: cfg.VPNFingerprint,
		client: &http.Client{
			Timeout: 15 * time.Second,
			Jar:     jar,
		},
	}
}

func (p *ThreeXUIVPNProvider) Name() string {
	return "3xui"
}

func (p *ThreeXUIVPNProvider) IssueKey(userID uint, endAt time.Time) (IssuedVPNKey, error) {
	if err := p.validateConfig(); err != nil {
		return IssuedVPNKey{}, err
	}

	if err := p.login(); err != nil {
		return IssuedVPNKey{}, err
	}

	clientID, err := p.generateUUID()
	if err != nil {
		return IssuedVPNKey{}, fmt.Errorf("generate client uuid: %w", err)
	}

	email := fmt.Sprintf("vpn-user-%d", userID)
	settings := threeXUIClientSettings{
		Clients: []threeXUIClient{
			{
				ID:         clientID,
				Flow:       p.flow,
				Email:      email,
				LimitIP:    0,
				TotalGB:    0,
				ExpiryTime: endAt.UnixMilli(),
				Enable:     true,
				TGID:       "",
				SubID:      "",
				Comment:    fmt.Sprintf("bot-user-%d", userID),
				Reset:      0,
			},
		},
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return IssuedVPNKey{}, fmt.Errorf("marshal 3x-ui client settings: %w", err)
	}

	requestBody, err := json.Marshal(threeXUIAddClientRequest{
		ID:       p.inboundID,
		Settings: string(settingsJSON),
	})
	if err != nil {
		return IssuedVPNKey{}, fmt.Errorf("marshal 3x-ui add client request: %w", err)
	}

	resp, err := p.client.Post(p.panelURL+"/panel/api/inbounds/addClient", "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return IssuedVPNKey{}, fmt.Errorf("3x-ui add client request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return IssuedVPNKey{}, fmt.Errorf("read 3x-ui add client response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return IssuedVPNKey{}, fmt.Errorf("3x-ui add client failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var addResp threeXUIBaseResponse
	if len(body) > 0 {
		if err := json.Unmarshal(body, &addResp); err == nil && !addResp.Success {
			return IssuedVPNKey{}, fmt.Errorf("3x-ui add client failed: %s", strings.TrimSpace(addResp.Msg))
		}
	}

	accessURL := p.buildAccessURL(clientID, email)
	configJSON, err := json.Marshal(map[string]string{
		"provider":    "3xui",
		"client_id":   clientID,
		"email":       email,
		"access_url":  accessURL,
		"server_name": p.serverName,
		"public_host": p.publicHost,
		"public_port": p.publicPort,
		"flow":        p.flow,
	})
	if err != nil {
		return IssuedVPNKey{}, fmt.Errorf("marshal 3x-ui client config: %w", err)
	}

	return IssuedVPNKey{
		ExternalClientID: clientID,
		KeyName:          email,
		AccessURL:        accessURL,
		ConfigJSON:       string(configJSON),
	}, nil
}

func (p *ThreeXUIVPNProvider) DeactivateKey(externalClientID string) error {
	if strings.TrimSpace(externalClientID) == "" {
		return nil
	}
	if err := p.validateConfig(); err != nil {
		return err
	}
	if err := p.login(); err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s/panel/api/inbounds/%d/delClient/%s", p.panelURL, p.inboundID, url.PathEscape(externalClientID))
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		return fmt.Errorf("build 3x-ui delete client request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("3x-ui delete client request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read 3x-ui delete client response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("3x-ui delete client failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var deleteResp threeXUIBaseResponse
	if len(body) > 0 {
		if err := json.Unmarshal(body, &deleteResp); err == nil && !deleteResp.Success {
			return fmt.Errorf("3x-ui delete client failed: %s", strings.TrimSpace(deleteResp.Msg))
		}
	}

	return nil
}

func (p *ThreeXUIVPNProvider) validateConfig() error {
	switch {
	case p.panelURL == "":
		return fmt.Errorf("3x-ui panel endpoint is not configured")
	case p.username == "":
		return fmt.Errorf("3x-ui panel username is not configured")
	case p.password == "":
		return fmt.Errorf("3x-ui panel password is not configured")
	case p.inboundID <= 0:
		return fmt.Errorf("3x-ui inbound id is not configured")
	case p.publicHost == "":
		return fmt.Errorf("3x-ui public host is not configured")
	case p.serverName == "":
		return fmt.Errorf("3x-ui reality server name is not configured")
	case p.publicKey == "":
		return fmt.Errorf("3x-ui reality public key is not configured")
	case p.shortID == "":
		return fmt.Errorf("3x-ui reality short id is not configured")
	default:
		return nil
	}
}

func (p *ThreeXUIVPNProvider) login() error {
	requestBody, err := json.Marshal(threeXUILoginRequest{
		Username: p.username,
		Password: p.password,
	})
	if err != nil {
		return fmt.Errorf("marshal 3x-ui login request: %w", err)
	}

	resp, err := p.client.Post(p.panelURL+"/login", "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("3x-ui login request failed: %w", err)
	}
	body, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if readErr != nil {
		return fmt.Errorf("read 3x-ui login response: %w", readErr)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var loginResp threeXUIBaseResponse
		if len(body) == 0 || json.Unmarshal(body, &loginResp) != nil || loginResp.Success || strings.TrimSpace(loginResp.Msg) == "" {
			return nil
		}
	}

	form := url.Values{}
	form.Set("username", p.username)
	form.Set("password", p.password)
	resp, err = p.client.Post(p.panelURL+"/login", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("3x-ui login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read 3x-ui login response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("3x-ui login failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var loginResp threeXUIBaseResponse
	if len(body) > 0 && json.Unmarshal(body, &loginResp) == nil && !loginResp.Success && strings.TrimSpace(loginResp.Msg) != "" {
		return fmt.Errorf("3x-ui login failed: %s", strings.TrimSpace(loginResp.Msg))
	}

	return nil
}

func (p *ThreeXUIVPNProvider) buildAccessURL(clientID string, email string) string {
	query := url.Values{}
	query.Set("security", "reality")
	query.Set("sni", p.serverName)
	query.Set("fp", p.fingerprint)
	query.Set("pbk", p.publicKey)
	query.Set("sid", p.shortID)
	query.Set("type", "tcp")
	query.Set("flow", p.flow)
	query.Set("encryption", "none")

	return fmt.Sprintf(
		"vless://%s@%s:%s?%s#%s",
		clientID,
		p.publicHost,
		p.publicPort,
		query.Encode(),
		url.QueryEscape(email),
	)
}

func (p *ThreeXUIVPNProvider) generateUUID() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		raw[0:4],
		raw[4:6],
		raw[6:8],
		raw[8:10],
		raw[10:16],
	), nil
}
