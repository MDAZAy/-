package models

import "time"

type Plan struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"size:255;not null"`
	Price        float64   `gorm:"not null"`
	DurationDays int       `gorm:"not null"`
	Description  string    `gorm:"type:text"`
	IsActive     bool      `gorm:"default:true"`
	CreatedAt    time.Time `gorm:"not null"`
}
