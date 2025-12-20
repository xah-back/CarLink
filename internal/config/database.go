package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetUpDatabaseConnection(logger *slog.Logger) *gorm.DB {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", "error", err)
		panic(err)
	}

	dbUrl := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbUrl,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		panic(err)
	}

	logger.Info("Successfully connected to the database")
	return db
}
