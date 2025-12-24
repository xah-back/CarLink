package transports

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/services"
)

func RegisterRoutes(
	routes *gin.Engine,
	logger *slog.Logger,
	userService services.UserService,
	carService services.CarService,
	tripService services.TripService,
	bookingService services.BookingService,
	reviewService services.ReviewService,
) {
	userHandler := NewUserHandler(userService, logger)
	carHandler := NewCarHandler(carService, logger)
	tripHandler := NewTripHandler(tripService, logger)
	bookingHandler := NewBookingHandler(bookingService, logger)
	reviewHandler := NewReviewHandler(reviewService, logger)

	userHandler.RegisterRoutes(routes)
	carHandler.RegisterRoutes(routes)
	tripHandler.RegisterRoutes(routes)
	bookingHandler.RegisterRoutes(routes)
	reviewHandler.RegisterRoutes(routes)
}
