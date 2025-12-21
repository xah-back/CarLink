package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/config"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
	"github.com/mutsaevz/team-5-ambitious/internal/services"
	"github.com/mutsaevz/team-5-ambitious/internal/transports"
)

func main() {
	// инициализация логгера (tmp внутри logging)
	logger := config.InitLogger()

	r := gin.Default()

	db := config.SetUpDatabaseConnection(logger)
	if db == nil {
		logger.Error("database is nil")
		return
	}

	if err := db.AutoMigrate(&models.User{}, &models.Car{}); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	logger.Info("migrations completed")

	userRepo := repository.NewUserRepository(db, logger)
	carRepo := repository.NewCarRepository(db, logger)

	userService := services.NewUserService(userRepo, logger)
	carService := services.NewCarService(carRepo, userRepo, logger)

	transports.RegisterRoutes(r, logger, userService, carService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("application started successfully")

	if err := r.Run(":" + port); err != nil {
		logger.Error("ошибка запуска сервера", slog.Any("error", err))
	}
}
