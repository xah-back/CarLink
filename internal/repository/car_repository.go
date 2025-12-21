package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"gorm.io/gorm"
)

type CarRepository interface {
	Create(car *models.Car) error

	// List() ([]models.Car, error)

	GetByID(id uint) (*models.Car, error)

	// Update(car *models.Car) (*models.Car, error)

	// Delete(id uint) error
}

type gormCarRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewCarRepository(db *gorm.DB, logger *slog.Logger) CarRepository {
	return &gormCarRepository{
		db:     db,
		logger: logger,
	}
}

func (r *gormCarRepository) Create(car *models.Car) error {

	err := r.db.Create(&car).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *gormCarRepository) GetByID(id uint) (*models.Car, error) {
	var car models.Car

	if err := r.db.First(&car, id).Error; err != nil {
		return nil, err
	}

	return &car, nil
}
