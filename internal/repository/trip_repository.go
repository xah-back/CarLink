package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"gorm.io/gorm"
)

type TripRepository interface {
	Create(trip *models.Trip) error

	List(filter models.TripFilter) ([]models.Trip, error)

	GetByID(id uint) (*models.Trip, error)

	Update(trip *models.Trip) error

	Delete(id uint) error

	WithDB(db *gorm.DB) TripRepository
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

func (r *gormTripRepository) List(filter models.TripFilter) ([]models.Trip, error) {
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

func (r *gormTripRepository) GetByID(id uint) (*models.Trip, error) {
	var trip models.Trip

	if err := r.db.First(&trip, id).Error; err != nil {
		return nil, err
	}

	return &trip, nil
}

func (r *gormTripRepository) Update(trip *models.Trip) error {
	op := "repository.booking.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("trip_id", uint64(trip.ID)),
	)

	return r.db.
		Model(&models.Booking{}).
		Where("id = ?", trip.ID).
		Updates(trip).
		Error
}

func (r *gormTripRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Trip{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormTripRepository) WithDB(db *gorm.DB) TripRepository {
	return &gormTripRepository{
		db:     db,
		logger: r.logger,
	}
}
