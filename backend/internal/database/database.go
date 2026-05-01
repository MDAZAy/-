package database

import (
	"database/sql"
	"errors"
	"path/filepath"
	"time"

	"autoservice/backend/internal/config"
	"autoservice/backend/internal/models"

	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormlogger "gorm.io/gorm/logger"
)

func Connect(cfg config.Config) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, sqlDB, nil
}

func ApplyMigrations(sqlDB *sql.DB, migrationsDir string) error {
	driver, err := migratemysql.WithInstance(sqlDB, &migratemysql.Config{})
	if err != nil {
		return err
	}

	absDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return err
	}

	instance, err := migrate.NewWithDatabaseInstance("file://"+filepath.ToSlash(absDir), "mysql", driver)
	if err != nil {
		return err
	}

	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func Seed(db *gorm.DB) error {
	adminRole := models.Role{Name: "admin"}
	customerRole := models.Role{Name: "customer"}

	if err := upsertRole(db, &adminRole); err != nil {
		return err
	}
	if err := upsertRole(db, &customerRole); err != nil {
		return err
	}
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "customer").First(&customerRole).Error; err != nil {
		return err
	}

	statuses := []models.AppointmentStatus{
		{Name: "Scheduled", Code: "scheduled"},
		{Name: "Completed", Code: "completed"},
		{Name: "Cancelled", Code: "cancelled"},
	}
	for i := range statuses {
		if err := upsertByCode(db, &statuses[i]); err != nil {
			return err
		}
	}

	categories := []models.ServiceCategory{
		{Name: "Maintenance", Description: "Regular maintenance and consumables replacement"},
		{Name: "Diagnostics", Description: "Computer diagnostics and fault finding"},
	}
	for i := range categories {
		if err := upsertCategory(db, &categories[i]); err != nil {
			return err
		}
	}

	var maintenance models.ServiceCategory
	if err := db.Where("name = ?", "Maintenance").First(&maintenance).Error; err != nil {
		return err
	}
	var diagnostics models.ServiceCategory
	if err := db.Where("name = ?", "Diagnostics").First(&diagnostics).Error; err != nil {
		return err
	}

	services := []models.Service{
		{CategoryID: maintenance.ID, Name: "Oil change", Description: "Engine oil and filter replacement", DurationMinutes: 60, Price: 89.90, IsActive: true},
		{CategoryID: maintenance.ID, Name: "Brake inspection", Description: "Brake pads and discs diagnostics", DurationMinutes: 45, Price: 59.90, IsActive: true},
		{CategoryID: diagnostics.ID, Name: "Computer diagnostics", Description: "Electronic systems diagnostics", DurationMinutes: 30, Price: 49.90, IsActive: true},
	}
	for i := range services {
		if err := upsertService(db, &services[i]); err != nil {
			return err
		}
	}

	for weekday := 0; weekday <= 6; weekday++ {
		entry := models.WorkingHour{
			Weekday:   weekday,
			StartTime: "09:00",
			EndTime:   "18:00",
			IsWorking: weekday != 0,
		}
		if err := upsertWorkingHour(db, &entry); err != nil {
			return err
		}
	}

	mechanics := []models.Mechanic{
		{FullName: "Alex Petrov", Phone: "+358400000001", Email: "alex.petrov@example.com", IsActive: true},
		{FullName: "Mika Virtanen", Phone: "+358400000002", Email: "mika.virtanen@example.com", IsActive: true},
	}
	for i := range mechanics {
		if err := upsertMechanic(db, &mechanics[i]); err != nil {
			return err
		}
	}

	adminHash, err := bcrypt.GenerateFromPassword([]byte("Admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userHash, err := bcrypt.GenerateFromPassword([]byte("User123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Email:        "admin@example.com",
		PasswordHash: string(adminHash),
		FullName:     "Service Administrator",
		Phone:        "+358400000010",
		RoleID:       adminRole.ID,
		IsActive:     true,
	}
	user := models.User{
		Email:        "user@example.com",
		PasswordHash: string(userHash),
		FullName:     "Default Client",
		Phone:        "+358400000020",
		RoleID:       customerRole.ID,
		IsActive:     true,
	}

	if err := upsertUser(db, &admin); err != nil {
		return err
	}
	if err := upsertUser(db, &user); err != nil {
		return err
	}
	if err := db.Where("email = ?", "admin@example.com").First(&admin).Error; err != nil {
		return err
	}
	if err := db.Where("email = ?", "user@example.com").First(&user).Error; err != nil {
		return err
	}

	vehicle := models.Vehicle{
		UserID:      user.ID,
		Make:        "Toyota",
		Model:       "Corolla",
		Year:        2021,
		PlateNumber: "AUTO-001",
		Color:       "White",
		VIN:         "JTDBR32E720012345",
	}
	if err := upsertVehicle(db, &vehicle); err != nil {
		return err
	}

	setting := models.AppSetting{Key: "display_timezone", Value: "Europe/Helsinki"}
	if err := upsertSetting(db, &setting); err != nil {
		return err
	}

	return nil
}

func upsertRole(db *gorm.DB, model *models.Role) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]any{
			"is_deleted": false,
			"deleted_at": nil,
			"updated_at": time.Now().UTC(),
		}),
	}).Create(model).Error
}

func upsertCategory(db *gorm.DB, model *models.ServiceCategory) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]any{
			"description": model.Description,
			"is_deleted":  false,
			"deleted_at":  nil,
			"updated_at":  time.Now().UTC(),
		}),
	}).Create(model).Error
}

func upsertByCode(db *gorm.DB, model *models.AppointmentStatus) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
	}).Create(model).Error
}

func upsertService(db *gorm.DB, model *models.Service) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]any{
			"category_id":      model.CategoryID,
			"description":      model.Description,
			"duration_minutes": model.DurationMinutes,
			"price":            model.Price,
			"is_active":        model.IsActive,
			"is_deleted":       false,
			"deleted_at":       nil,
			"updated_at":       time.Now().UTC(),
		}),
	}).Create(model).Error
}

func upsertWorkingHour(db *gorm.DB, model *models.WorkingHour) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "weekday"}},
		DoUpdates: clause.Assignments(map[string]any{
			"start_time": model.StartTime,
			"end_time":   model.EndTime,
			"is_working": model.IsWorking,
			"updated_at": time.Now().UTC(),
		}),
	}).Create(model).Error
}

func upsertMechanic(db *gorm.DB, model *models.Mechanic) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.Assignments(map[string]any{
			"full_name":  model.FullName,
			"phone":      model.Phone,
			"is_active":  model.IsActive,
			"is_deleted": false,
			"deleted_at": nil,
			"updated_at": time.Now().UTC(),
		}),
	}).Create(model).Error
}

func upsertUser(db *gorm.DB, model *models.User) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.Assignments(map[string]any{
			"password_hash": model.PasswordHash,
			"full_name":     model.FullName,
			"phone":         model.Phone,
			"role_id":       model.RoleID,
			"is_active":     model.IsActive,
			"is_deleted":    false,
			"deleted_at":    nil,
			"updated_at":    time.Now().UTC(),
		}),
	}).Create(model).Error
}

func upsertVehicle(db *gorm.DB, model *models.Vehicle) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "plate_number"}},
		DoUpdates: clause.Assignments(map[string]any{
			"user_id":    model.UserID,
			"make":       model.Make,
			"model":      model.Model,
			"year":       model.Year,
			"color":      model.Color,
			"vin":        model.VIN,
			"is_deleted": false,
			"deleted_at": nil,
			"updated_at": time.Now().UTC(),
		}),
	}).Create(model).Error
}

func upsertSetting(db *gorm.DB, model *models.AppSetting) error {
	if model.ID == "" {
		model.ID = uuid.NewString()
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(model).Error
}
