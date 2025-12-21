package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/config"
	"github.com/mutsaevz/team-5-ambitious/internal/logging"
)

func main() {
	// инициализация логгера (tmp внутри logging)
	logger := logging.New()

	r := gin.Default()

	db := config.SetUpDatabaseConnection(logger)
	if db == nil {
		logger.Error("database is nil")
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("application started successfully")

	if err := r.Run(":" + port); err != nil {
		logger.Error("ошибка запуска сервера", slog.Any("error", err))
	}
}
