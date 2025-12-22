package services

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
	"gorm.io/gorm"
)

type BookingService interface {
	Create(req *models.BookingCreateRequest) (*models.Booking, error)

	List() ([]models.Booking, error)

	GetByID(id uint) (*models.Booking, error)

	Update(id uint, req *models.BookingUpdateRequest) (*models.Booking, error)

	Delete(id uint) error
}

type bookingService struct {
	repo   repository.BookingRepository
	db     *gorm.DB
	logger *slog.Logger
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	db *gorm.DB,
	logger *slog.Logger,
) BookingService {
	return &bookingService{
		repo:   bookingRepo,
		db:     db,
		logger: logger,
	}
}

func (s *bookingService) Create(req *models.BookingCreateRequest) (*models.Booking, error) {

	op := "service.booking.create"

	s.logger.Debug("repo call",
		slog.String("op", op),
		slog.Uint64("trip_id", uint64(req.TripID)),
		slog.Uint64("passenger_id", uint64(req.PassengerID)),
	)

	booking := &models.Booking{
		TripID:      req.TripID,
		PassengerID: req.PassengerID,
	}
	if err := s.repo.Create(booking); err != nil {
		s.logger.Error("repo error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("booking created", slog.String("op", op), slog.Uint64("id", uint64(booking.ID)))
	return booking, nil
}

func (s *bookingService) List() ([]models.Booking, error) {

	op := "service.booking.list"

	s.logger.Debug("repo call", slog.String("op", op))

	bookings, err := s.repo.List()
	if err != nil {
		s.logger.Error("repo error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("bookings listed", slog.String("op", op), slog.Int("count", len(bookings)))
	return bookings, nil
}

func (s *bookingService) GetByID(id uint) (*models.Booking, error) {
	op := "service.booking.GetByID"

	s.logger.Debug("repo call", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))

	booking, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("repo error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	return booking, nil
}

func (s *bookingService) Update(id uint, req *models.BookingUpdateRequest) (*models.Booking, error) {
	op := "service.booking.Update"

	s.logger.Debug("repo call", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))

	booking, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("repo error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	if req.BookingStatus != nil {
		booking.BookingStatus = *req.BookingStatus
	}

	if err := s.repo.Update(booking); err != nil {
		s.logger.Error("repo error", slog.String("op", op), slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("booking updated", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))
	return booking, err
}

func (s *bookingService) Delete(id uint) error {

	op := "service.booking.Delete"
	s.logger.Debug("repo call", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("repo error", slog.String("op", op), slog.Any("error", err))
		return err
	}
	s.logger.Info("booking deleted", slog.String("op", op), slog.Uint64("booking_id", uint64(id)))
	return nil
}
