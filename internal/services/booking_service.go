package services

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/constants"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
	"gorm.io/gorm"
)

type BookingService interface {
	Create(tripID uint, passengerID uint) (*models.Booking, error)

	List() ([]models.Booking, error)

	Approve(bookingID uint, driverID uint) error

	Rejected(bookingID uint, driverID uint) error

	GetByID(id uint) (*models.Booking, error)

	Update(id uint, req *models.BookingUpdateRequest) (*models.Booking, error)

	Delete(id uint) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	tripRepo    repository.TripRepository
	db          *gorm.DB
	logger      *slog.Logger
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	tripRepo repository.TripRepository,
	db *gorm.DB,
	logger *slog.Logger,
) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		tripRepo:    tripRepo,
		db:          db,
		logger:      logger,
	}
}

func (s *bookingService) Create(
	tripID uint,
	passengerID uint,
) (*models.Booking, error) {

	trip, err := s.tripRepo.GetByID(tripID)
	if err != nil {
		return nil, fmt.Errorf("trip not found: %w", err)
	}

	if constants.TripStatus(trip.TripStatus) != constants.TripPublished {
		return nil, errors.New("trip is not available")
	}

	exists, err := s.bookingRepo.Exists(tripID, passengerID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("request already sent")
	}

	booking := &models.Booking{
		TripID:        tripID,
		PassengerID:   passengerID,
		BookingStatus: constants.BookingPending,
	}

	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) Approve(bookingID, driverID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		bookingRepo := s.bookingRepo.WithDB(tx)
		tripRepo := s.tripRepo.WithDB(tx)

		booking, err := bookingRepo.GetByID(bookingID)
		if err != nil {
			return err
		}

		if booking.BookingStatus != constants.BookingPending {
			return errors.New("booking is not pending")
		}

		trip, err := tripRepo.GetByID(booking.TripID)
		if err != nil {
			return err
		}

		if trip.DriverID != driverID {
			return errors.New("forbidden")
		}

		if trip.AvailableSeats <= 0 {
			return errors.New("no available seats")
		}

		// Водитель одобряет
		trip.AvailableSeats--
		booking.BookingStatus = constants.BookingApproved

		err = tripRepo.Update(trip)
		if err != nil {
			return err
		}

		err = bookingRepo.Update(booking)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *bookingService) Rejected(bookingID uint, driverID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		bookingRepo := s.bookingRepo.WithDB(tx)
		tripRepo := s.tripRepo.WithDB(tx)

		booking, err := bookingRepo.GetByID(bookingID)
		if err != nil {
			return fmt.Errorf("booking not found: %w", err)
		}

		trip, err := tripRepo.GetByID(booking.TripID)
		if err != nil {
			return fmt.Errorf("trip not found: %w", err)
		}

		if trip.DriverID != driverID {
			return errors.New("forbidden")
		}

		if booking.BookingStatus != constants.BookingPending {
			return errors.New("booking is not pending")
		}

		booking.BookingStatus = constants.BookingRejected

		if err := bookingRepo.Update(booking); err != nil {
			return err
		}

		return nil
	})
}

func (s *bookingService) List() ([]models.Booking, error) {

	op := "service.booking.list"

	s.logger.Debug(" call", slog.String("op", op))

	bookings, err := s.bookingRepo.List()
	if err != nil {
		s.logger.Error(" error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("bookings listed", slog.String("op", op), slog.Int("count", len(bookings)))
	return bookings, nil
}

func (s *bookingService) GetByID(id uint) (*models.Booking, error) {
	op := "service.booking.GetByID"

	s.logger.Debug(" call", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))

	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		s.logger.Error(" error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	return booking, nil
}

func (s *bookingService) Update(id uint, req *models.BookingUpdateRequest) (*models.Booking, error) {
	op := "service.booking.Update"

	s.logger.Debug(" call", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))

	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		s.logger.Error(" error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	if req.BookingStatus != nil {
		booking.BookingStatus = *req.BookingStatus
	}

	if err := s.bookingRepo.Update(booking); err != nil {
		s.logger.Error(" error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("booking updated", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))
	return booking, err
}

func (s *bookingService) Delete(id uint) error {

	op := "service.booking.Delete"
	s.logger.Debug(" call", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))

	if err := s.bookingRepo.Delete(id); err != nil {
		s.logger.Error(" error", slog.String("op", op), slog.Any("error", err))
		return err
	}
	s.logger.Info("booking deleted", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))
	return nil
}
