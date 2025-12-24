package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
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
		api.POST("/", h.Create)
		api.GET("/", h.List)
		api.GET("/:id", h.GetByID)
		api.GET("/:id", h.GetByID)
		api.GET("/trip/:trip_id/pending", h.GetAllPendingBookingsByTripID)
		api.PATCH("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
	}
}

func (h *BookingHandler) Create(ctx *gin.Context) {

	h.logger.Info("handler called",
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.FullPath()),
	)

	var input models.BookingCreateRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("invalid JSON",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	booking, err := h.service.Create(&input)
	if err != nil {
		h.logger.Error("error adding booking",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err.Error()),
		)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("booking created successfully")
	ctx.JSON(http.StatusCreated, booking)
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

func (h *BookingHandler) GetAllPendingBookingsByTripID(ctx *gin.Context) {

	h.logger.Info("handler called",
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.FullPath()),
	)
	tripIDParam := ctx.Param("trip_id")

	tripID, err := strconv.ParseUint(tripIDParam, 10, 64)

	if err != nil {
		h.logger.Warn("invalid trip ID parameter",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID parameter"})
		return
	}

	bookings, err := h.service.GetAllPendingBookingsByTripID(uint(tripID))

	if err != nil {
		h.logger.Error("error getting pending bookings by trip ID",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err.Error()),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("pending bookings retrieved successfully")
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

func (h *BookingHandler) Update(ctx *gin.Context) {

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

	var input models.BookingUpdateRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("invalid JSON",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	booking, err := h.service.Update(uint(id), &input)

	if err != nil {
		h.logger.Error("error updating booking",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err.Error()),
		)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("booking updated successfully")
	ctx.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) Delete(ctx *gin.Context) {

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

	if err := h.service.Delete(uint(id)); err != nil {
		h.logger.Error("error deleting booking",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err.Error()),
		)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.Info("booking deleted successfully")
	ctx.JSON(http.StatusNoContent, nil)
}
