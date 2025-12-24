package repository

import (
	"errors"
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/dto"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"gorm.io/gorm"
)

type TripRepository interface {
	Create(trip *models.Trip) error

	List(filter dto.TripFilter) ([]models.Trip, error)

	GetByID(id uint) (*models.Trip, error)

	Update(trip *models.Trip) error

	Delete(id uint) error

	WithDB(db *gorm.DB) TripRepository

	UpdateAvgRating(tripID uint, avg float64) error

	IsPassenger(tripID, userID uint) (bool, error)
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

	if err := r.db.Create(trip).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormTripRepository) List(filter dto.TripFilter) ([]models.Trip, error) {
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

	if filter.StartTime != nil {
		query = query.Where("start_time >= ?", *filter.StartTime)
	}

	if filter.TripStatus != nil {
		query = query.Where("trip_status = ?", *filter.TripStatus)
	}

	page := filter.Page
	pageSize := filter.PageSize

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *gormTripRepository) GetByID(id uint) (*models.Trip, error) {
	var trip models.Trip

	if err := r.db.First(&trip, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &trip, nil
}

func (r *gormTripRepository) Update(trip *models.Trip) error {
	op := "repository.trip.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("trip_id", uint64(trip.ID)),
	)

	return r.db.
		Model(&models.Trip{}).
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

func (r *gormTripRepository) UpdateAvgRating(tripID uint, avg float64) error {
	op := "repository.trip.update_avg_rating"
	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("trip_id", uint64(tripID)),
		slog.Float64("avg_rating", avg),
	)
	if err := r.db.Model(&models.Trip{}).Where("id = ?", tripID).Update("avg_rating", avg).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}
	return nil
}

func (r *gormTripRepository) IsPassenger(tripID, userID uint) (bool, error) {
	var count int64

	err := r.db.Model(&models.Booking{}).
		Where("trip_id = ? AND passenger_id = ?", tripID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
