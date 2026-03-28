package database

import (
	// "fmt"
	"log"

	"techbrein-project-management/internal/config"
	"techbrein-project-management/internal/domain/entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SeedAdmin(db *gorm.DB, cfg config.Config) error {
	var count int64
	if err := db.Model(&entities.User{}).Where("email = ?", cfg.DefaultAdminEmail).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.DefaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := entities.User{
		Name:         cfg.DefaultAdminName,
		Email:        cfg.DefaultAdminEmail,
		PasswordHash: string(hash),
		Role:         entities.RoleAdmin,
	}

	log.Printf("seeding default admin user %s", cfg.DefaultAdminEmail)
	return db.Create(&admin).Error
}

func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func DsnExample() string {
	return "postgres://postgres:postgres@localhost:5432/techbrein_pm?sslmode=disable"
}
