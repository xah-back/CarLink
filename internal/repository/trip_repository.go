package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"gorm.io/gorm"
)

type TripRepository interface {
	Create(trip *models.Trip) error

	GetByQueryParameters(filter models.TripFilter) ([]models.Trip, error)

	// List() ([]models.Trip, error)

	// GetByID(id uint) (*models.Trip, error)

	// Update(trip *models.Trip) (*models.Trip, error)

	// Delete(id uint) error
}

type gormTripRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTripRepository(db *gorm.DB, logger *slog.Logger) TripRepository {
	return &gormTripRepository{
		db:     db,
		logger: logger,
	}
}

func (r *gormTripRepository) Create(trip *models.Trip) error {

	if err := r.db.Create(&trip).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormTripRepository) GetByQueryParameters(filter models.TripFilter) ([]models.Trip, error) {
	var list []models.Trip

	query := r.db.Model(&models.Trip{}).
		Where("available_seats > 0")

	if filter.FromCity != nil {
		query = query.Where("from_city = ?", *filter.FromCity)
	}

	if filter.ToCity != nil {
		query = query.Where("to_city = ?", *filter.ToCity)
	}

	if filter.AvailableSeats != nil {
		query = query.Where("available_seats >= ?", *filter.AvailableSeats)
	}

	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}
