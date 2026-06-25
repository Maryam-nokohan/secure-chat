package postgres

import (
	"fmt"

	"github.com/maryam-nokohan/secure-chat/internal/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *configs.Config) (*gorm.DB, error) {

	defaultDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBPort,
	)

	tempDB, err := gorm.Open(
		postgres.Open(defaultDSN),
		&gorm.Config{},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect postgres database: %w",
			err,
		)
	}

	createDBQuery := fmt.Sprintf(
		"CREATE DATABASE %s",
		cfg.DBName,
	)

	err = tempDB.Exec(createDBQuery).Error
	if err != nil {
		fmt.Println("database probably already exists:", err)
	}

	db, err := gorm.Open(
		postgres.Open(cfg.DSN),
		&gorm.Config{},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect app database: %w",
			err,
		)
	}

	return db, nil
}