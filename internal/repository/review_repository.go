package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error

	List() ([]models.Review, error)

	GetByID(id uint) (*models.Review, error)

	Update(review *models.Review) (*models.Review, error)

	Delete(id uint) error

	ExistsByTripAndUser(tripID, userID uint) (bool, error)

	GetAvgRatingByTrip(tripID uint) (float64, error)
}

type gormReviewRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewReviewRepository(db *gorm.DB, logger *slog.Logger) ReviewRepository {
	return &gormReviewRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *gormReviewRepository) Create(review *models.Review) error {
	op := "repository.review.create"
	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("rating", uint64(review.Rating)),
		slog.String("text", review.Text),
	)
	if err := r.DB.Create(review).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}
	return nil
}

func (r *gormReviewRepository) List() ([]models.Review, error) {

	op := "repository.review.list"
	r.logger.Debug("db call", slog.String("op", op))
	var reviews []models.Review

	if err := r.DB.Find(&reviews).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	r.logger.Debug("db response", slog.String("op", op), slog.Int("count", len(reviews)))
	return reviews, nil
}

func (r *gormReviewRepository) GetByID(id uint) (*models.Review, error) {

	op := "repository.review.get_by_id"
	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("id", uint64(id)),
	)
	var review models.Review

	if err := r.DB.Where("id = ?", id).First(&review).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	return &review, nil
}

func (r *gormReviewRepository) Update(review *models.Review) (*models.Review, error) {

	op := "repository.review.update"
	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("id", uint64(review.ID)),
	)
	if err := r.DB.Model(&models.Review{}).Where("id = ?", review.ID).Updates(review).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	return review, nil
}

func (r *gormReviewRepository) Delete(id uint) error {

	op := "repository.review.delete"
	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("id", uint64(id)),
	)
	result := r.DB.Delete(&models.Review{}, id)
	if result.Error != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", result.Error))
		return result.Error
	}
	return nil
}

func (r *gormReviewRepository) ExistsByTripAndUser(tripID, userID uint) (bool, error) {

	op := "repository.review.exists_by_trip_and_user"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("trip_id", uint64(tripID)),
		slog.Uint64("user_id", uint64(userID)),
	)
	var count int64

	if err := r.DB.Model(&models.Review{}).
		Where("trip_id = ? AND author_id = ?", tripID, userID).
		Count(&count).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return false, err
	}
	return count > 0, nil

}

func (r *gormReviewRepository) GetAvgRatingByTrip(tripID uint) (float64, error) {

	op := "repository.review.get_avg_rating_by_trip"
	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("trip_id", uint64(tripID)),
	)
	var avgRating float64

	if err := r.DB.Model(&models.Review{}).
		Where("trip_id = ?", tripID).
		Select("AVG(rating)").
		Scan(&avgRating).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return 0, err
	}
	return avgRating, nil
}
