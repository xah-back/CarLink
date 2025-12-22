package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"gorm.io/gorm"
)

type TripRepository interface {
	Create(trip *models.Trip) error

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
