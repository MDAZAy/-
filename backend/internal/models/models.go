package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string     `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;autoUpdateTime" json:"updated_at"`
	IsDeleted bool       `gorm:"not null;default:false;index" json:"is_deleted"`
	DeletedAt *time.Time `gorm:"default:null" json:"deleted_at,omitempty"`
}

func (b *BaseModel) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.NewString()
	}
	return nil
}

type Role struct {
	BaseModel
	Name string `gorm:"size:100;not null;uniqueIndex" json:"name"`
}

type User struct {
	BaseModel
	Email        string `gorm:"size:255;not null;uniqueIndex:users_email" json:"email"`
	PasswordHash string `gorm:"type:text;not null" json:"-"`
	FullName     string `gorm:"size:255;not null" json:"full_name"`
	Phone        string `gorm:"size:50;not null" json:"phone"`
	RoleID       string `gorm:"type:char(36);not null;index" json:"role_id"`
	Role         Role   `gorm:"foreignKey:RoleID" json:"role"`
	IsActive     bool   `gorm:"not null;default:true" json:"is_active"`
	Vehicles     []Vehicle
}

type Vehicle struct {
	BaseModel
	UserID      string `gorm:"type:char(36);not null;index:vehicles_user" json:"user_id"`
	Make        string `gorm:"size:120;not null" json:"make"`
	Model       string `gorm:"size:120;not null" json:"model"`
	Year        int    `gorm:"not null" json:"year"`
	PlateNumber string `gorm:"size:50;not null;uniqueIndex" json:"plate_number"`
	Color       string `gorm:"size:80" json:"color"`
	VIN         string `gorm:"size:50" json:"vin"`
}

type ServiceCategory struct {
	BaseModel
	Name        string `gorm:"size:150;not null;uniqueIndex" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

type Service struct {
	BaseModel
	CategoryID      string          `gorm:"type:char(36);not null;index" json:"category_id"`
	Category        ServiceCategory `gorm:"foreignKey:CategoryID" json:"category"`
	Name            string          `gorm:"size:150;not null" json:"name"`
	Description     string          `gorm:"type:text" json:"description"`
	DurationMinutes int             `gorm:"not null" json:"duration_minutes"`
	Price           float64         `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive        bool            `gorm:"not null;default:true" json:"is_active"`
}

type AppointmentStatus struct {
	BaseModel
	Name string `gorm:"size:100;not null" json:"name"`
	Code string `gorm:"size:50;not null;uniqueIndex" json:"code"`
}

type WorkingHour struct {
	BaseModel
	Weekday   int    `gorm:"not null;uniqueIndex" json:"weekday"`
	StartTime string `gorm:"size:5;not null" json:"start_time"`
	EndTime   string `gorm:"size:5;not null" json:"end_time"`
	IsWorking bool   `gorm:"not null;default:true" json:"is_working"`
}

type Holiday struct {
	BaseModel
	HolidayDate time.Time `gorm:"type:date;not null;uniqueIndex" json:"holiday_date"`
	Name        string    `gorm:"size:255;not null" json:"name"`
}

type Mechanic struct {
	BaseModel
	FullName string `gorm:"size:255;not null" json:"full_name"`
	Phone    string `gorm:"size:50;not null" json:"phone"`
	Email    string `gorm:"size:255;not null;uniqueIndex" json:"email"`
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`
}

type Appointment struct {
	BaseModel
	UserID             string            `gorm:"type:char(36);not null;index:appointments_user" json:"user_id"`
	User               User              `gorm:"foreignKey:UserID" json:"user"`
	VehicleID          string            `gorm:"type:char(36);not null" json:"vehicle_id"`
	Vehicle            Vehicle           `gorm:"foreignKey:VehicleID" json:"vehicle"`
	ServiceID          string            `gorm:"type:char(36);not null" json:"service_id"`
	Service            Service           `gorm:"foreignKey:ServiceID" json:"service"`
	StatusID           string            `gorm:"type:char(36);not null;index:appointments_status" json:"status_id"`
	Status             AppointmentStatus `gorm:"foreignKey:StatusID" json:"status"`
	MechanicID         string            `gorm:"type:char(36);not null" json:"mechanic_id"`
	Mechanic           Mechanic          `gorm:"foreignKey:MechanicID" json:"mechanic"`
	StartTime          time.Time         `gorm:"not null;index:appointments_time,priority:1" json:"start_time"`
	EndTime            time.Time         `gorm:"not null;index:appointments_time,priority:2" json:"end_time"`
	ConfirmationNumber string            `gorm:"size:50;not null;uniqueIndex" json:"confirmation_number"`
	IdempotencyKey     string            `gorm:"size:255;not null;uniqueIndex:idx_appointments_user_idempotency" json:"-"`
	Notes              string            `gorm:"type:text" json:"notes"`
}

type AppointmentHistory struct {
	BaseModel
	AppointmentID   string  `gorm:"type:char(36);not null;index" json:"appointment_id"`
	StatusID        string  `gorm:"type:char(36);not null" json:"status_id"`
	ChangedByUserID *string `gorm:"type:char(36)" json:"changed_by_user_id,omitempty"`
	Comment         string  `gorm:"type:text" json:"comment"`
}

func (AppointmentHistory) TableName() string {
	return "appointment_history"
}

type Notification struct {
	BaseModel
	UserID        string  `gorm:"type:char(36);not null;index" json:"user_id"`
	AppointmentID *string `gorm:"type:char(36)" json:"appointment_id,omitempty"`
	Type          string  `gorm:"size:100;not null" json:"type"`
	Message       string  `gorm:"type:text;not null" json:"message"`
	IsRead        bool    `gorm:"not null;default:false" json:"is_read"`
}

type Payment struct {
	BaseModel
	AppointmentID string     `gorm:"type:char(36);not null;index" json:"appointment_id"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency      string     `gorm:"size:10;not null;default:EUR" json:"currency"`
	Status        string     `gorm:"size:50;not null" json:"status"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}

type AppointmentFile struct {
	BaseModel
	AppointmentID string `gorm:"type:char(36);not null;index" json:"appointment_id"`
	FileName      string `gorm:"size:255;not null" json:"file_name"`
	FileURL       string `gorm:"type:text;not null" json:"file_url"`
}

type AuditLog struct {
	BaseModel
	UserID      *string `gorm:"type:char(36);index" json:"user_id,omitempty"`
	Action      string  `gorm:"size:150;not null" json:"action"`
	Entity      string  `gorm:"size:150;not null" json:"entity"`
	EntityID    string  `gorm:"size:64;not null" json:"entity_id"`
	IPAddress   string  `gorm:"size:64;not null" json:"ip_address"`
	Metadata    string  `gorm:"type:text;not null" json:"metadata"`
	Description string  `gorm:"type:text;not null" json:"description"`
}

type RefreshToken struct {
	BaseModel
	UserID     string     `gorm:"type:char(36);not null;index" json:"user_id"`
	TokenHash  string     `gorm:"size:128;not null;uniqueIndex" json:"-"`
	ExpiresAt  time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	UserAgent  string     `gorm:"size:255;not null" json:"user_agent"`
	IPAddress  string     `gorm:"size:64;not null" json:"ip_address"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

type AppSetting struct {
	BaseModel
	Key   string `gorm:"size:120;not null;uniqueIndex" json:"key"`
	Value string `gorm:"type:text;not null" json:"value"`
}
