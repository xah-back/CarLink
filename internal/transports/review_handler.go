package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/services"
)

type ReviewHandler struct {
	service services.ReviewService
	logger  *slog.Logger
}

func NewReviewHandler(service services.ReviewService, logger *slog.Logger) *ReviewHandler {
	return &ReviewHandler{
		service: service,
		logger:  logger,
	}
}

func (h ReviewHandler) RegisterRoutes(ctx *gin.Engine) {
	api := ctx.Group("")
	{
		api.POST("/trips/:id/reviews", h.Create)
	}
}

func (h *ReviewHandler) Create(ctx *gin.Context) {

	tripIDParam := ctx.Param("id")
	tripID, err := strconv.ParseUint(tripIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		h.logger.Error("invalid trip id",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		return
	}

	var req models.ReviewCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logger.Error("invalid request body",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		return
	}

	review, err := h.service.Create(uint(tripID), &req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logger.Error("error creating review",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		return
	}
	h.logger.Info("review created successfully",
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.FullPath()),
	)
	ctx.JSON(http.StatusCreated, review)

}
