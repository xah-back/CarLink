package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/config"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
	"github.com/mutsaevz/team-5-ambitious/internal/services"
	"github.com/mutsaevz/team-5-ambitious/internal/transports"
	"github.com/redis/go-redis/v9"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// инициализация логгера (tmp внутри logging)
	logger := config.InitLogger()

	r := gin.New()
	r.Use(gin.Recovery())

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	db := config.SetUpDatabaseConnection(logger)
	if db == nil {
		logger.Error("database is nil")
		return
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Car{},
		&models.Trip{},
		&models.Booking{},
		&models.Review{}); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	logger.Info("migrations completed")

	userRepo := repository.NewUserRepository(db, logger)
	carRepo := repository.NewCarRepository(db, logger)
	tripRepo := repository.NewTripRepository(db, logger)

	tripStatusWorker := services.NewTripStatusWorker(
		tripRepo,
		logger,
		time.Minute,
	)

	tripStatusWorker.Start(ctx)

	bookingRepo := repository.NewBookingRepository(db, logger)
	reviewRepo := repository.NewReviewRepository(db, logger)

	userService := services.NewUserService(userRepo, logger)
	carService := services.NewCarService(carRepo, userRepo, logger)
	tripService := services.NewTripService(tripRepo, userRepo, carRepo, logger)
	bookingService := services.NewBookingService(bookingRepo, tripRepo, db, logger)
	reviewService := services.NewReviewService(reviewRepo, tripRepo, db, rdb, logger)

	transports.RegisterRoutes(
		r, logger,
		userService,
		carService,
		tripService,
		bookingService,
		reviewService,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("application started successfully")

	if err := r.Run(":" + port); err != nil {
		logger.Error("ошибка запуска сервера", slog.Any("error", err))
	}
}
