package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"autoservice/backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewManager(cfg config.Config) *Manager {
	return &Manager{
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessTTL:     cfg.AccessTTL,
		refreshTTL:    cfg.RefreshTTL,
	}
}

func (m *Manager) GenerateTokenPair(userID, role, email string) (string, string, time.Time, error) {
	accessToken, err := m.generateToken(userID, role, email, "access", m.accessTTL, m.accessSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshToken, err := m.generateToken(userID, role, email, "refresh", m.refreshTTL, m.refreshSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, time.Now().UTC().Add(m.refreshTTL), nil
}

func (m *Manager) ParseAccessToken(token string) (*Claims, error) {
	return m.parse(token, m.accessSecret, "access")
}

func (m *Manager) ParseRefreshToken(token string) (*Claims, error) {
	return m.parse(token, m.refreshSecret, "refresh")
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (m *Manager) generateToken(userID, role, email, tokenType string, ttl time.Duration, secret []byte) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (m *Manager) parse(raw string, secret []byte, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(raw, &Claims{}, func(_ *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.Type != expectedType {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
