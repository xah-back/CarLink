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

	{
		api.POST("/:id", h.Create)
		api.GET("/:id", h.GetByID)
	}
}

func (h *CarHandler) Create(ctx *gin.Context) {
	var input models.CarCreateRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	car, err := h.service.Create(uint(id), input)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "car create error"})
		return
	}

	ctx.JSON(http.StatusOK, car)
}

func (h *CarHandler) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	car, err := h.service.GetByID(uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "record not found"})
		return
	}

	ctx.JSON(http.StatusOK, car)
}
