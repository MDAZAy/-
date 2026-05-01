package services

import (
	"net/http"
	"strings"
	"time"

	"autoservice/backend/internal/auth"
	"autoservice/backend/internal/config"
	"autoservice/backend/internal/dto"
	"autoservice/backend/internal/models"
	"autoservice/backend/internal/repositories"
	"autoservice/backend/internal/validators"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RequestMeta struct {
	IPAddress string
	UserAgent string
}

type AuthService struct {
	repo       *repositories.Repository
	jwtManager *auth.Manager
	cfg        config.Config
}

func NewAuthService(repo *repositories.Repository, jwtManager *auth.Manager, cfg config.Config) *AuthService {
	return &AuthService{repo: repo, jwtManager: jwtManager, cfg: cfg}
}

func (s *AuthService) Register(req dto.RegisterRequest, meta RequestMeta) (*dto.AuthResponse, *AppError) {
	if err := validators.ValidateRegister(req); err != nil {
		return nil, NewError(http.StatusBadRequest, "validation_failed", err.Error())
	}

	if _, err := s.repo.FindUserByEmail(strings.TrimSpace(strings.ToLower(req.Email))); err == nil {
		return nil, NewError(http.StatusConflict, "email_taken", "user with this email already exists")
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, NewError(http.StatusInternalServerError, "user_lookup_failed", "failed to inspect user")
	}

	role, err := s.repo.FindRoleByName("customer")
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "role_not_found", "customer role is not configured")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "password_hash_failed", "failed to protect password")
	}

	user := models.User{
		Email:        strings.TrimSpace(strings.ToLower(req.Email)),
		PasswordHash: string(hash),
		FullName:     strings.TrimSpace(req.FullName),
		Phone:        strings.TrimSpace(req.Phone),
		RoleID:       role.ID,
		IsActive:     true,
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, NewError(http.StatusInternalServerError, "user_create_failed", "failed to create user")
	}

	user.Role = *role
	return s.issueTokens(&user, meta, "register")
}

func (s *AuthService) Login(req dto.LoginRequest, meta RequestMeta) (*dto.AuthResponse, *AppError) {
	if err := validators.ValidateLogin(req); err != nil {
		return nil, NewError(http.StatusBadRequest, "validation_failed", err.Error())
	}

	user, err := s.repo.FindUserByEmail(strings.TrimSpace(strings.ToLower(req.Email)))
	if err != nil {
		return nil, NewError(http.StatusUnauthorized, "invalid_credentials", "invalid credentials")
	}
	if !user.IsActive {
		return nil, NewError(http.StatusForbidden, "user_inactive", "user is inactive")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, NewError(http.StatusUnauthorized, "invalid_credentials", "invalid credentials")
	}

	return s.issueTokens(user, meta, "login")
}

func (s *AuthService) Refresh(req dto.RefreshRequest, meta RequestMeta) (*dto.AuthResponse, *AppError) {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return nil, NewError(http.StatusBadRequest, "validation_failed", "refresh_token is required")
	}

	claims, err := s.jwtManager.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, NewError(http.StatusUnauthorized, "invalid_refresh_token", "refresh token is invalid")
	}

	hashed := auth.HashToken(req.RefreshToken)
	stored, err := s.repo.FindRefreshToken(hashed)
	if err != nil {
		return nil, NewError(http.StatusUnauthorized, "refresh_token_not_found", "refresh token is invalid")
	}

	if stored.RevokedAt != nil || stored.ExpiresAt.Before(time.Now().UTC()) {
		return nil, NewError(http.StatusUnauthorized, "refresh_token_expired", "refresh token is expired")
	}

	if err := s.repo.RevokeRefreshToken(hashed); err != nil {
		return nil, NewError(http.StatusInternalServerError, "refresh_token_revoke_failed", "failed to rotate refresh token")
	}

	user, err := s.repo.FindUserByID(claims.UserID)
	if err != nil {
		return nil, NewError(http.StatusUnauthorized, "user_not_found", "user not found")
	}

	return s.issueTokens(user, meta, "refresh")
}

func (s *AuthService) Logout(req dto.LogoutRequest) *AppError {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return NewError(http.StatusBadRequest, "validation_failed", "refresh_token is required")
	}
	if err := s.repo.RevokeRefreshToken(auth.HashToken(req.RefreshToken)); err != nil {
		return NewError(http.StatusInternalServerError, "logout_failed", "failed to revoke refresh token")
	}
	return nil
}

func (s *AuthService) issueTokens(user *models.User, meta RequestMeta, action string) (*dto.AuthResponse, *AppError) {
	accessToken, refreshToken, refreshExpiry, err := s.jwtManager.GenerateTokenPair(user.ID, user.Role.Name, user.Email)
	if err != nil {
		return nil, NewError(http.StatusInternalServerError, "token_issue_failed", "failed to issue tokens")
	}

	tokenRecord := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: auth.HashToken(refreshToken),
		ExpiresAt: refreshExpiry,
		UserAgent: limitText(meta.UserAgent, 255),
		IPAddress: limitText(meta.IPAddress, 64),
	}
	if err := s.repo.CreateRefreshToken(&tokenRecord); err != nil {
		return nil, NewError(http.StatusInternalServerError, "refresh_token_store_failed", "failed to store refresh token")
	}

	description := "User authenticated successfully"
	if action == "register" {
		description = "User account created"
	}
	_ = s.repo.CreateAuditLog(&models.AuditLog{
		UserID:      &user.ID,
		Action:      action,
		Entity:      "auth",
		EntityID:    user.ID,
		IPAddress:   limitText(meta.IPAddress, 64),
		Metadata:    `{"email":"` + user.Email + `"}`,
		Description: description,
	})

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.AuthUser{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			Role:     user.Role.Name,
		},
	}, nil
}

func limitText(value string, size int) string {
	value = strings.TrimSpace(value)
	if len(value) <= size {
		return value
	}
	return value[:size]
}
