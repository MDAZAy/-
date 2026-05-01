package repositories

import (
	"time"

	"autoservice/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) FindRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ? AND is_deleted = ?", name, false).First(&role).Error
	return &role, err
}

func (r *Repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	return &user, err
}

func (r *Repository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("id = ? AND is_deleted = ?", userID, false).First(&user).Error
	return &user, err
}

func (r *Repository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *Repository) FindRefreshToken(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Where("token_hash = ? AND is_deleted = ?", hash, false).First(&token).Error
	return &token, err
}

func (r *Repository) RevokeRefreshToken(hash string) error {
	now := time.Now().UTC()
	return r.db.Model(&models.RefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", hash).
		Updates(map[string]any{"revoked_at": &now, "updated_at": now}).
		Error
}

func (r *Repository) ListActiveCategories() ([]models.ServiceCategory, error) {
	var categories []models.ServiceCategory
	err := r.db.Where("is_deleted = ?", false).Order("name").Find(&categories).Error
	return categories, err
}

func (r *Repository) ListActiveServices() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Preload("Category").
		Where("is_active = ? AND is_deleted = ?", true, false).
		Order("name").
		Find(&services).Error
	return services, err
}

func (r *Repository) CreateVehicle(vehicle *models.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *Repository) ListUserVehicles(userID string) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := r.db.Where("user_id = ? AND is_deleted = ?", userID, false).Order("created_at DESC").Find(&vehicles).Error
	return vehicles, err
}

func (r *Repository) CreateAuditLog(entry *models.AuditLog) error {
	return r.db.Create(entry).Error
}

func (r *Repository) CreateAuditLogTx(tx *gorm.DB, entry *models.AuditLog) error {
	return tx.Create(entry).Error
}

func (r *Repository) CreateNotificationTx(tx *gorm.DB, notification *models.Notification) error {
	return tx.Create(notification).Error
}

func (r *Repository) CreateAppointmentHistoryTx(tx *gorm.DB, history *models.AppointmentHistory) error {
	return tx.Create(history).Error
}

func (r *Repository) FindStatusByCodeTx(tx *gorm.DB, code string) (*models.AppointmentStatus, error) {
	var status models.AppointmentStatus
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("code = ? AND is_deleted = ?", code, false).
		First(&status).Error
	return &status, err
}

func (r *Repository) FindVehicleForUserTx(tx *gorm.DB, userID, vehicleID string) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND user_id = ? AND is_deleted = ?", vehicleID, userID, false).
		First(&vehicle).Error
	return &vehicle, err
}

func (r *Repository) FindServiceTx(tx *gorm.DB, serviceID string) (*models.Service, error) {
	var service models.Service
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Category").
		Where("id = ? AND is_active = ? AND is_deleted = ?", serviceID, true, false).
		First(&service).Error
	return &service, err
}

func (r *Repository) FindWorkingHourTx(tx *gorm.DB, weekday int) (*models.WorkingHour, error) {
	var item models.WorkingHour
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("weekday = ? AND is_deleted = ?", weekday, false).
		First(&item).Error
	return &item, err
}

func (r *Repository) IsHolidayTx(tx *gorm.DB, day time.Time) (bool, error) {
	var count int64
	err := tx.Model(&models.Holiday{}).
		Where("holiday_date = ? AND is_deleted = ?", day.Format("2006-01-02"), false).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) FindAppointmentByIdempotencyTx(tx *gorm.DB, userID, key string) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.preloadAppointment(tx.Clauses(clause.Locking{Strength: "UPDATE"})).
		Where("user_id = ? AND idempotency_key = ? AND is_deleted = ?", userID, key, false).
		First(&appointment).Error
	return &appointment, err
}

func (r *Repository) FindFreeMechanicTx(tx *gorm.DB, start, end time.Time) (*models.Mechanic, error) {
	subquery := tx.Model(&models.Appointment{}).
		Select("1").
		Where("appointments.mechanic_id = mechanics.id").
		Where("appointments.is_deleted = ?", false).
		Where("appointments.start_time < ? AND appointments.end_time > ?", end, start)

	var mechanic models.Mechanic
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("is_active = ? AND is_deleted = ?", true, false).
		Where("NOT EXISTS (?)", subquery).
		Order("full_name").
		First(&mechanic).Error
	return &mechanic, err
}

func (r *Repository) CreateAppointmentTx(tx *gorm.DB, appointment *models.Appointment) error {
	return tx.Create(appointment).Error
}

func (r *Repository) LoadAppointment(appointment *models.Appointment) error {
	return r.preloadAppointment(r.db).First(appointment, "id = ?", appointment.ID).Error
}

func (r *Repository) ListUserAppointments(userID string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.preloadAppointment(r.db).
		Where("user_id = ? AND is_deleted = ?", userID, false).
		Order("start_time DESC").
		Find(&appointments).Error
	return appointments, err
}

func (r *Repository) ListAllAppointments() ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.preloadAppointment(r.db).
		Where("is_deleted = ?", false).
		Order("start_time DESC").
		Find(&appointments).Error
	return appointments, err
}

func (r *Repository) ListDayAppointments(dayStart, dayEnd time.Time) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.preloadAppointment(r.db).
		Where("is_deleted = ?", false).
		Where("start_time < ? AND end_time > ?", dayEnd, dayStart).
		Find(&appointments).Error
	return appointments, err
}

func (r *Repository) CountActiveMechanics() (int64, error) {
	var count int64
	err := r.db.Model(&models.Mechanic{}).Where("is_active = ? AND is_deleted = ?", true, false).Count(&count).Error
	return count, err
}

func (r *Repository) GetService(serviceID string) (*models.Service, error) {
	var service models.Service
	err := r.db.Preload("Category").
		Where("id = ? AND is_active = ? AND is_deleted = ?", serviceID, true, false).
		First(&service).Error
	return &service, err
}

func (r *Repository) GetDashboardCounts() (map[string]int64, error) {
	result := map[string]int64{}
	var count int64
	if err := r.db.Model(&models.User{}).Where("is_deleted = ?", false).Count(&count).Error; err != nil {
		return nil, err
	}
	result["users"] = count
	if err := r.db.Model(&models.Vehicle{}).Where("is_deleted = ?", false).Count(&count).Error; err != nil {
		return nil, err
	}
	result["vehicles"] = count
	if err := r.db.Model(&models.Appointment{}).Where("is_deleted = ?", false).Count(&count).Error; err != nil {
		return nil, err
	}
	result["appointments"] = count
	if err := r.db.Model(&models.Mechanic{}).Where("is_deleted = ?", false).Count(&count).Error; err != nil {
		return nil, err
	}
	result["mechanics"] = count
	return result, nil
}

func (r *Repository) preloadAppointment(query *gorm.DB) *gorm.DB {
	return query.Preload("Status").
		Preload("Service.Category").
		Preload("Vehicle").
		Preload("Mechanic").
		Preload("User")
}
