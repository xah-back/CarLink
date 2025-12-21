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
) {
	userHandler := NewUserHandler(userService, logger)
	carHandler := NewCarHandler(carService, logger)

	userHandler.RegisterRoutes(routes)
	carHandler.RegisterRoutes(routes)
}
