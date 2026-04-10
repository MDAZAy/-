package storage

import (
	"context"
	"time"
)

// Repository describes the storage operations used by the agent backend.
type Repository interface {
	Save(ctx context.Context, input SaveHistoryInput) (*DBHistory, error)
	GetRecentSuccess(ctx context.Context, limit int) ([]DBHistory, error)
	GetStats(ctx context.Context, filter StatsFilter) (*Stats, error)
}

// StatsFilter limits aggregate queries by creation time.
type StatsFilter struct {
	From *time.Time
	To   *time.Time
}
