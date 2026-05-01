package middleware

import (
	"net/http"
	"sync"
	"time"

	"autoservice/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type limiterEntry struct {
	count   int
	resetAt time.Time
}

type inMemoryLimiter struct {
	mu     sync.Mutex
	items  map[string]limiterEntry
	limit  int
	window time.Duration
}

func NewRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := &inMemoryLimiter{
		items:  make(map[string]limiterEntry),
		limit:  limit,
		window: window,
	}

	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.FullPath()
		now := time.Now()

		limiter.mu.Lock()
		entry, ok := limiter.items[key]
		if !ok || now.After(entry.resetAt) {
			entry = limiterEntry{count: 0, resetAt: now.Add(limiter.window)}
		}
		entry.count++
		limiter.items[key] = entry
		limiter.mu.Unlock()

		if entry.count > limiter.limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.Envelope{Success: false, Error: "rate limit exceeded", Code: "rate_limited"})
			return
		}

		c.Next()
	}
}
