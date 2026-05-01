package middleware

import (
	"net/http"
	"strings"

	"autoservice/backend/internal/auth"
	"autoservice/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "user_role"
)

func AuthRequired(jwtManager *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Envelope{Success: false, Error: "missing bearer token", Code: "unauthorized"})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := jwtManager.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Envelope{Success: false, Error: "invalid access token", Code: "unauthorized"})
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentRole, _ := c.Get(ContextRole)
		if currentRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Envelope{Success: false, Error: "access denied", Code: "forbidden"})
			return
		}
		c.Next()
	}
}
