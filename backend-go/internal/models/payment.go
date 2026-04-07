package models

import "time"

type Payment struct {
	ID                uint      `gorm:"primaryKey"`
	UserID            uint      `gorm:"index;not null"`
	PlanID            uint      `gorm:"index;not null"`
	Amount            float64   `gorm:"not null"`
	Currency          string    `gorm:"size:16;not null"`
	Status            string    `gorm:"size:64;index;not null"`
	Provider          string    `gorm:"size:64;not null"`
	ExternalPaymentID string    `gorm:"size:255;uniqueIndex;not null"`
	PaymentURL        string    `gorm:"type:text"`
	RawResponse       string    `gorm:"type:text"`
	CreatedAt         time.Time `gorm:"not null"`
}
