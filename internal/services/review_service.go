package services

import (
	"errors"
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/constants"
	"github.com/mutsaevz/team-5-ambitious/internal/dto"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrTripNotCompleted     = errors.New("trip not completed")
	ErrUserNotPassenger     = errors.New("user is not a passenger in this trip")
	ErrReviewAlreadyPresent = errors.New("review already exists for this user and trip")
)

type ReviewService interface {
	Create(tripID uint, req *dto.ReviewCreateRequest) (*models.Review, error)

	List(filter models.Page) ([]models.Review, error)

	// GetByID(id uint) (*models.Review, error)

	// Update(id uint, req *models.Review) (*models.Review, error)

	// Delete(id uint) error
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	tripRepo   repository.TripRepository
	logger     *slog.Logger
	db         *gorm.DB
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	tripRepo repository.TripRepository,
	db *gorm.DB,
	logger *slog.Logger,
) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		tripRepo:   tripRepo,
		logger:     logger,
		db:         db,
	}
}

func (s *reviewService) Create(tripID uint, req *dto.ReviewCreateRequest) (*models.Review, error) {
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

		isPassenger, err := tr.IsPassenger(tripID, req.AuthorID)
		if err != nil {
			return err
		}

		if !isPassenger {
			s.logger.Error("user is not passenger",
				slog.String("op", op),
				slog.Uint64("userID", uint64(req.AuthorID)),
				slog.Uint64("tripID", uint64(tripID)),
			)
			return ErrUserNotPassenger
		}

		exists, err := rr.ExistsByTripAndUser(tripID, req.AuthorID)
		if err != nil {
			s.logger.Error("error checking review existence", slog.String("op", op), slog.Any("error", err))
			return err
		}
		if exists {
			s.logger.Error("review already exists for this user and trip",
				slog.String("op", op), slog.Uint64("userID", uint64(req.AuthorID)), slog.Uint64("tripID", uint64(tripID)))
			return ErrReviewAlreadyPresent
		}

		review := &models.Review{
			AuthorID: req.AuthorID,
			TripID:   tripID,
			Rating:   req.Rating,
			Text:     req.Text,
		}

		if err := rr.Create(review); err != nil {
			s.logger.Error("error creating review", slog.String("op", op), slog.Any("error", err))
			return err
		}

		avgRating, err := rr.GetAvgRatingByTrip(tripID)
		if err != nil {
			s.logger.Error("error calculating average rating", slog.String("op", op), slog.Any("error", err))
			return err
		}

		if err := tr.UpdateAvgRating(tripID, avgRating); err != nil {
			s.logger.Error("error updating trip average rating", slog.String("op", op), slog.Any("error", err))
			return err
		}

		created = review
		return nil
	})

	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *reviewService) List(filter models.Page) ([]models.Review, error) {

	op := "service.review.list"

	s.logger.Debug(" call", slog.String("op", op))

	reviews, err := s.reviewRepo.List(filter)
	if err != nil {
		s.logger.Error(" error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	s.logger.Info("reviews listed", slog.String("op", op), slog.Int("count", len(reviews)))
	return reviews, nil
}
