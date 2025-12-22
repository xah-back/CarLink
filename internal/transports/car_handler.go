package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/services"
)

type CarHandler struct {
	service services.CarService
	logger  *slog.Logger
}

func NewCarHandler(service services.CarService, logger *slog.Logger) *CarHandler {
	return &CarHandler{
		service: service,
		logger:  logger,
	}
}

func (h *CarHandler) RegisterRoutes(ctx *gin.Engine) {
	api := ctx.Group("/cars")

	api.POST("/:id", h.Create)
	api.GET("/", h.List)
	api.GET("/owner/:id", h.GetByOwner)
	api.GET("/:id", h.GetByID)
	api.PUT("/:id", h.Update)
	api.DELETE("/:id", h.Delete)
}

// POST /cars/:id
func (h *CarHandler) Create(ctx *gin.Context) {
	var input models.CarCreateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid JSON for car creation", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid owner ID for car creation", slog.String("id", idStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	car, err := h.service.Create(uint(id), input)
	if err != nil {
		h.logger.Error("Failed to create car", slog.Uint64("owner_id", id), slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "car create error"})
		return
	}

	h.logger.Info("Car created successfully", slog.Uint64("car_id", car.ID), slog.Uint64("owner_id", id))
	ctx.JSON(http.StatusOK, car)
}

// GET /cars
func (h *CarHandler) List(ctx *gin.Context) {
	cars, err := h.service.List()
	if err != nil {
		h.logger.Error("Failed to list cars", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cars"})
		return
	}

	h.logger.Info("List of cars retrieved", slog.Int("count", len(cars)))
	ctx.JSON(http.StatusOK, cars)
}

// GET /cars/owner/:id
func (h *CarHandler) GetByOwner(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid owner ID for get by owner", slog.String("id", idStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	car, err := h.service.GetByOwner(uint(id))
	if err != nil {
		h.logger.Warn("Car not found for owner", slog.Uint64("owner_id", id))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	h.logger.Info("Car retrieved by owner", slog.Uint64("owner_id", id), slog.Uint64("car_id", car.ID))
	ctx.JSON(http.StatusOK, car)
}

// GET /cars/:id
func (h *CarHandler) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid car ID for get by ID", slog.String("id", idStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	car, err := h.service.GetByID(uint(id))
	if err != nil {
		h.logger.Warn("Car not found by ID", slog.Uint64("car_id", id))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	h.logger.Info("Car retrieved by ID", slog.Uint64("car_id", car.ID))
	ctx.JSON(http.StatusOK, car)
}

// PUT /cars/:id
func (h *CarHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid car ID for update", slog.String("id", idStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var input models.CarUpdateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid JSON for car update", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	car, err := h.service.Update(uint(id), input)
	if err != nil {
		h.logger.Error("Failed to update car", slog.Uint64("car_id", id), slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	h.logger.Info("Car updated successfully", slog.Uint64("car_id", car.ID))
	ctx.JSON(http.StatusOK, car)
}

// DELETE /cars/:id
func (h *CarHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid car ID for delete", slog.String("id", idStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		h.logger.Error("Failed to delete car", slog.Uint64("car_id", id), slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	h.logger.Info("Car deleted successfully", slog.Uint64("car_id", id))
	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
