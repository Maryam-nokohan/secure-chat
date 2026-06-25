package migrations

import (
	"os"

	"github.com/maryam-nokohan/secure-chat/pkg"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {

	sqlBytes, err := os.ReadFile(
		"internal/adapters/secondary/postgres/migrations/schema.sql",
	)
	if err != nil {
		return err
	}

	if err := db.Exec(string(sqlBytes)).Error; err != nil {
		return err
	}

	pkg.LogInfo("database migrations completed")

	return nil
}