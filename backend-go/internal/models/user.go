package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey"`
	TelegramID int64     `gorm:"uniqueIndex;not null"`
	Username   string    `gorm:"size:255"`
	FullName   string    `gorm:"size:255"`
	IsAdmin    bool      `gorm:"default:false"`
	IsBlocked  bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"not null"`
}
