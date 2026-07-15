package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/maryam-nokohan/secure-chat/pkg"
)


func Load() (*Config, error) {
	pkg.LogInfo("Loading configuration...")

	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			pkg.LogInfo(".env file not found, relying on process environment (expected in containers)")
		} else {
			return nil, fmt.Errorf("failed to load .env: %w", err)
		}
	}
	cfg := &Config{
	DBUser:      os.Getenv("DB_USER"),
	DBPassword:  os.Getenv("DB_PASSWORD"),
	DBHost:      os.Getenv("DB_HOST"),
	DBPort:      os.Getenv("DB_PORT"),
	DBName:      os.Getenv("DB_NAME"),
	RedisURL:   os.Getenv("REDIS_URL"),
	NatsURL:     os.Getenv("NATS_URL"),
	JWTSecret:   os.Getenv("JWT_SECRET"),
	CSRFSecrete: os.Getenv("CSRF_SECRET"),
	CacheEncryptionKey: os.Getenv("CACHE_ENCRYPTION_KEY"),
}

	cfg.DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	return cfg, nil
}