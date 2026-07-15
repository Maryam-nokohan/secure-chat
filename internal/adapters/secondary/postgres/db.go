package postgres

import (
	"fmt"

	"github.com/maryam-nokohan/secure-chat/internal/configs"
	"github.com/maryam-nokohan/secure-chat/pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *configs.Config) (*gorm.DB, error) {
	pkg.LogInfo("Connecting to database...")
	defaultDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort,
	)

	tempDB, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres database: %w", err)
	}

	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)"
	if err := tempDB.Raw(checkQuery, cfg.DBName).Scan(&exists).Error; err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
		if err := tempDB.Exec(createDBQuery).Error; err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		pkg.LogInfo("Database " + cfg.DBName + " created.")
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect app database: %w", err)
	}

	return db, nil
}