package repository

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *models.Booking) error

	List() ([]models.Booking, error)

	GetByID(id uint) (*models.Booking, error)

	Exists(tripID uint, passengerID uint) (bool, error)

	Update(booking *models.Booking) error

	Delete(id uint) error

	WithDB(db *gorm.DB) BookingRepository
}

type gormBookingRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewBookingRepository(db *gorm.DB, logger *slog.Logger) BookingRepository {
	return &gormBookingRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *gormBookingRepository) Create(booking *models.Booking) error {
	op := "repository.booking.create"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("trip_id", uint64(booking.TripID)),
		slog.Uint64("passenger_id", uint64(booking.PassengerID)),
	)

	if err := r.DB.Create(booking).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormBookingRepository) List() ([]models.Booking, error) {

	op := "repository.booking.list"

	r.logger.Debug("db call", slog.String("op", op))

	var bookings []models.Booking

	if err := r.DB.Find(&bookings).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}

	return bookings, nil

}

func (r *gormBookingRepository) GetByID(id uint) (*models.Booking, error) {

	op := "repository.booking.get_by_id"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("booking_id", uint64(id)),
	)

	var booking models.Booking

	if err := r.DB.Where("id = ?", id).First(&booking).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	return &booking, nil
}

func (r *gormBookingRepository) Exists(tripID uint, passengerID uint) (bool, error) {

	op := "repository.booking.exists"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("trip_id", uint64(tripID)),
		slog.Uint64("passenger_id", uint64(passengerID)),
	)

	var count int64

	if err := r.DB.Model(&models.Booking{}).
		Where("trip_id = ? AND passenger_id = ?", tripID, passengerID).
		Count(&count).Error; err != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", err))
		return false, err
	}

	exists := count > 0

	return exists, nil
}

func (r *gormBookingRepository) Update(booking *models.Booking) error {

	op := "repository.booking.update"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("booking_id", uint64(booking.ID)),
	)

	return r.DB.
		Model(&models.Booking{}).
		Where("id = ?", booking.ID).
		Update("booking_status", booking.BookingStatus).
		Error
}

func (r *gormBookingRepository) Delete(id uint) error {
	op := "repository.booking.delete"

	r.logger.Debug("db call",
		slog.String("op", op),
		slog.Uint64("booking_id", uint64(id)),
	)

	result := r.DB.Delete(&models.Booking{}, id)

	if result.Error != nil {
		r.logger.Error("db error", slog.String("op", op), slog.Any("error", result.Error))
		return result.Error
	}
	return nil
}

func (r *gormBookingRepository) WithDB(db *gorm.DB) BookingRepository {
	return &gormBookingRepository{
		DB:     db,
		logger: r.logger,
	}
}
