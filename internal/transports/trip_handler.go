package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/services"
)

type TripHandler struct {
	service services.TripService
	logger  *slog.Logger
}

func NewTripHandler(service services.TripService, logger *slog.Logger) *TripHandler {
	return &TripHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TripHandler) RegisterRoutes(ctx *gin.Engine) {
	api := ctx.Group("/trips")
	{
		api.POST("/:id", h.Create)
		api.GET("/", h.List)
	}
}

func (h *TripHandler) Create(ctx *gin.Context) {
	var req models.TripCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	trip, err := h.service.Create(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, trip)
}

func (h *TripHandler) List(ctx *gin.Context) {
	var filter models.TripFilter

	if from := ctx.Query("from_city"); from != "" {
		filter.FromCity = &from
	}

	if to := ctx.Query("to_city"); to != "" {
		filter.ToCity = &to
	}

	list, err := h.service.List(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, list)
}
