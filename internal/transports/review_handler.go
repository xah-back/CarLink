package transports

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mutsaevz/team-5-ambitious/internal/dto"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
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
		api.POST("/trips/:id/:author_id/reviews", h.Create)
		api.GET("/reviews", h.List)
		api.GET("/reviews/:id", h.GetByID)
		api.PUT("/reviews/:id/:author_id", h.Update)
		api.DELETE("/reviews/:id/:author_id", h.Delete)
	}
}

func (h *ReviewHandler) Create(ctx *gin.Context) {

	tripID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		h.logger.Error("invalid trip id",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		return
	}

	authorID, err := strconv.ParseUint(ctx.Param("author_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
		h.logger.Error("invalid author id",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		return
	}

	var req dto.ReviewCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logger.Error("invalid request body",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		return
	}

	review, err := h.service.Create(uint(tripID), uint(authorID), &req)

	if err != nil {

		if err == repository.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
			return
		}

		if err == services.ErrTripNotCompleted || err == services.ErrUserNotPassenger || err == services.ErrReviewAlreadyPresent {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error("error creating review",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	h.logger.Info("review created successfully",
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.FullPath()),
	)
	ctx.JSON(http.StatusCreated, review)

}

func (h *ReviewHandler) List(ctx *gin.Context) {

	var filter models.Page

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logger.Error("invalid query parameters",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		return
	}
	reviews, err := h.service.List(filter)
	if err != nil {
		h.logger.Error("error listing reviews",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	h.logger.Info("reviews listed successfully",
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.FullPath()),
	)
	ctx.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) GetByID(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		h.logger.Error("invalid review id",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		return
	}
	review, err := h.service.GetByID(uint(id))
	if err != nil {
		if err == repository.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
			return
		}
		h.logger.Error("error retrieving review",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	h.logger.Info("review retrieved successfully",
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.FullPath()),
	)
	ctx.JSON(http.StatusOK, review)
}

func (h *ReviewHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("invalid review id",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}
	authorID, err := strconv.ParseUint(ctx.Param("author_id"), 10, 64)
	if err != nil {
		h.logger.Error("invalid author id",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
		return
	}

	var req dto.ReviewUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.service.Update(uint(id), uint(authorID), &req)
	if err != nil {
		if err == repository.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, review)
}

func (h *ReviewHandler) Delete(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

	if err != nil {
		h.logger.Error("invalid review id",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
	}
	authorID, err := strconv.ParseUint(ctx.Param("author_id"), 10, 64)

	if err != nil {
		h.logger.Error("invalid author id",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
		return
	}

	if err := h.service.Delete(uint(id), uint(authorID)); err != nil {
		if err == repository.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"status": "deleted"})
}
