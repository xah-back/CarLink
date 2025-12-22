package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/services"
)

type BookingHandler struct {
	service services.BookingService
	logger  *slog.Logger
}

func NewBookingHandler(service services.BookingService, logger *slog.Logger) *BookingHandler {
	return &BookingHandler{
		service: service,
		logger:  logger,
	}
}

func (h BookingHandler) RegisterRoutes(ctx *gin.Engine) {
	api := ctx.Group("/bookings")
	{
		api.GET("/", h.List)
		api.GET("/:id", h.GetByID)
	}
}

func (h *BookingHandler) List(ctx *gin.Context) {

	h.logger.Info("handler called",
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.FullPath()),
	)

	bookings, err := h.service.List()

	if err != nil {
		h.logger.Error("error getting bookings",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err.Error()),
		)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("bookings retrieved successfully")
	ctx.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) GetByID(ctx *gin.Context) {

	h.logger.Info("handler called",
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.FullPath()),
	)
	idParam := ctx.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		h.logger.Warn("invalid ID parameter",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID parameter"})
		return
	}

	booking, err := h.service.GetByID(uint(id))

	if err != nil {
		h.logger.Error("error getting booking by ID",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err.Error()),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("booking retrieved successfully")
	ctx.JSON(http.StatusOK, booking)
}
