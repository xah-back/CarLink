package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/mutsaevz/team-5-ambitious/internal/constants"
	"github.com/mutsaevz/team-5-ambitious/internal/dto"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrTripNotCompleted     = errors.New("trip not completed")
	ErrUserNotPassenger     = errors.New("user is not a passenger in this trip")
	ErrReviewAlreadyPresent = errors.New("review already exists for this user and trip")
)

type ReviewService interface {
	Create(tripID, authorID uint, req *dto.ReviewCreateRequest) (*models.Review, error)

	List(filter models.Page) ([]dto.ReviewListItem, error)

	GetByID(id uint) (*models.Review, error)

	Update(id, authorID uint, req *dto.ReviewUpdateRequest) (*models.Review, error)

	Delete(id, authorID uint) error
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	tripRepo   repository.TripRepository
	logger     *slog.Logger
	redis      *redis.Client
	db         *gorm.DB
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	tripRepo repository.TripRepository,
	db *gorm.DB,
	redis *redis.Client,
	logger *slog.Logger,
) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		tripRepo:   tripRepo,
		logger:     logger,
		redis:      redis,
		db:         db,
	}
}

func (s *reviewService) Create(tripID, authorId uint, req *dto.ReviewCreateRequest) (*models.Review, error) {
	op := "service.review.create"

	var created *models.Review

	err := s.db.Transaction(func(tx *gorm.DB) error {
		tr := s.tripRepo.WithDB(tx)
		rr := s.reviewRepo.WithDB(tx)

		trip, err := tr.GetByID(tripID)
		if err != nil {
			return err
		}

		if trip.TripStatus != string(constants.TripCompleted) {
			s.logger.Error("cannot review a trip that is not completed",
				slog.String("op", op),
				slog.Uint64("tripID", uint64(tripID)),
			)
			return ErrTripNotCompleted
		}

		isPassenger, err := tr.IsPassenger(tripID, authorId)
		if err != nil {
			return err
		}

		if !isPassenger {
			s.logger.Error("user is not passenger",
				slog.String("op", op),
				slog.Uint64("userID", uint64(authorId)),
				slog.Uint64("tripID", uint64(tripID)),
			)
			return ErrUserNotPassenger
		}

		exists, err := rr.ExistsByTripAndUser(tripID, authorId)
		if err != nil {
			s.logger.Error("error checking review existence", slog.String("op", op), slog.Any("error", err))
			return err
		}
		if exists {
			s.logger.Error("review already exists for this user and trip",
				slog.String("op", op), slog.Uint64("userID", uint64(authorId)), slog.Uint64("tripID", uint64(tripID)))
			return ErrReviewAlreadyPresent
		}

		review := &models.Review{
			AuthorID: authorId,
			TripID:   tripID,
			Rating:   req.Rating,
			Text:     req.Text,
		}

		if err := rr.Create(review); err != nil {
			s.logger.Error("error creating review", slog.String("op", op), slog.Any("error", err))
			return err
		}

		if err := tr.UpdateAvgRatingFromReviews(tripID); err != nil {
			s.logger.Error("error updating trip average rating", slog.String("op", op), slog.Any("error", err))
			return err
		}

		created = review
		return nil
	})

	if err != nil {
		return nil, err
	}
	s.invalidateReviewListCache()
	return created, nil
}

func (s *reviewService) List(filter models.Page) ([]dto.ReviewListItem, error) {
	op := "service.review.list"

	// 🔥 кешируем ТОЛЬКО первую страницу БЕЗ фильтров
	useCache := filter.Page <= 1 &&
		filter.LastID == nil &&
		filter.TripID == nil &&
		filter.AuthorID == nil

	ctx := context.Background()
	var cacheKey string

	if useCache {
		cacheKey = buildReviewListCacheKey(filter)

		if data, err := s.redis.Get(ctx, cacheKey).Bytes(); err == nil {
			var items []dto.ReviewListItem
			if err := json.Unmarshal(data, &items); err == nil {
				s.logger.Debug("cache hit", slog.String("op", op))
				return items, nil
			}
		}
	}

	items, err := s.reviewRepo.List(filter)
	if err != nil {
		return nil, err
	}

	if useCache {
		data, _ := json.Marshal(items)
		ttl := 30*time.Second + time.Duration(rand.Intn(10))*time.Second
		_ = s.redis.Set(ctx, cacheKey, data, ttl).Err()
	}

	return items, nil
}

func (s *reviewService) GetByID(id uint) (*models.Review, error) {
	op := "service.review.getByID"
	s.logger.Debug("call", slog.String("op", op), slog.Uint64("id", uint64(id)))

	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		s.logger.Error("error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("review retrieved", slog.String("op", op), slog.Uint64("id", uint64(id)))
	return review, nil
}

func (s *reviewService) Update(id, authorID uint, req *dto.ReviewUpdateRequest) (*models.Review, error) {
	op := "service.review.update"

	var updated *models.Review

	err := s.db.Transaction(func(tx *gorm.DB) error {
		rr := s.reviewRepo.WithDB(tx)
		tr := s.tripRepo.WithDB(tx)

		review, err := rr.GetByID(id)
		if err != nil {
			return err
		}

		if review.AuthorID != authorID {
			return errors.New("permission denied")
		}

		if req.Text != nil {
			review.Text = *req.Text
		}
		if req.Rating != nil {
			review.Rating = *req.Rating
		}

		if _, err := rr.Update(review); err != nil {
			s.logger.Error("error updating review", slog.String("op", op), slog.Any("error", err))
			return err
		}

		if err := tr.UpdateAvgRatingFromReviews(review.TripID); err != nil {
			return err
		}

		updated = review
		return nil
	})

	if err != nil {
		return nil, err
	}
	s.invalidateReviewListCache()
	return updated, nil
}

func (s *reviewService) Delete(id, authorID uint) error {
	op := "service.review.delete"

	err := s.db.Transaction(func(tx *gorm.DB) error {
		rr := s.reviewRepo.WithDB(tx)
		tr := s.tripRepo.WithDB(tx)

		review, err := rr.GetByID(id)
		if err != nil {
			return err
		}

		if review.AuthorID != authorID {
			return errors.New("permission denied")
		}

		if err := rr.Delete(id); err != nil {
			s.logger.Error("error deleting review", slog.String("op", op), slog.Any("error", err))
			return err
		}

		return tr.UpdateAvgRatingFromReviews(review.TripID)
	})
	if err != nil {
		return err
	}
	s.invalidateReviewListCache()
	return nil

}

func buildReviewListCacheKey(filter models.Page) string {
	tripID := uint(0)
	if filter.TripID != nil {
		tripID = *filter.TripID
	}

	authorID := uint(0)
	if filter.AuthorID != nil {
		authorID = *filter.AuthorID
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return fmt.Sprintf(
		"reviews:list:trip=%d:author=%d:page=%d:size=%d",
		tripID,
		authorID,
		page,
		pageSize,
	)
}

func (s *reviewService) invalidateReviewListCache() {
	ctx := context.Background()

	pageSizes := []int{10, 20, 50, 100}
	for _, size := range pageSizes {
		key := fmt.Sprintf(
			"reviews:list:trip=0:author=0:page=1:size=%d",
			size,
		)
		_ = s.redis.Del(ctx, key).Err()
	}
}
