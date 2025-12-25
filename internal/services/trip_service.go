package services

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/constants"
	"github.com/mutsaevz/team-5-ambitious/internal/dto"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
)

type TripService interface {
	Create(driverID uint, req *dto.TripCreateRequest) (*models.Trip, error)

	List(filter dto.TripFilter) ([]models.Trip, error)

	GetByID(id uint) (*models.Trip, error)

	Update(id uint, req dto.TripUpdateRequest) (*models.Trip, error)

	Delete(id uint) error
}

type tripService struct {
	tripRepo repository.TripRepository
	userRepo repository.UserRepository
	carRepo  repository.CarRepository
	logger   *slog.Logger
}

func NewTripService(
	tripRepo repository.TripRepository,
	userRepo repository.UserRepository,
	carRepo repository.CarRepository,
	logger *slog.Logger) TripService {
	return &tripService{
		tripRepo: tripRepo,
		userRepo: userRepo,
		carRepo:  carRepo,
		logger:   logger,
	}
}

func (s *tripService) Create(id uint, req *dto.TripCreateRequest) (*models.Trip, error) {
	driver, err := s.userRepo.GetByID(id)

	if err != nil {
		return nil, err
	}

	car, err := s.carRepo.GetByOwner(id)

	if err != nil {
		return nil, err
	}

	var trip = models.Trip{
		DriverID:       driver.ID,
		CarID:          car.ID,
		FromCity:       req.FromCity,
		ToCity:         req.ToCity,
		StartTime:      req.StartTime,
		DurationMin:    req.DurationMin,
		TotalSeats:     car.Seats,
		AvailableSeats: req.AvailableSeats,
		Price:          req.Price,
		TripStatus:     string(constants.TripPublished),
		AvgRating:      0,
	}

	if err := s.tripRepo.Create(&trip); err != nil {
		return nil, err
	}

	return &trip, nil
}

func (s *tripService) List(filter dto.TripFilter) ([]models.Trip, error) {
	return s.tripRepo.List(filter)
}

func (s *tripService) GetByID(id uint) (*models.Trip, error) {
	trip, err := s.tripRepo.GetByID(id)
	if err != nil {
		s.logger.Error("trip not found",
			slog.Uint64("trip_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	return trip, nil
}

func (s *tripService) Update(id uint, req dto.TripUpdateRequest) (*models.Trip, error) {
	trip, err := s.tripRepo.GetByID(id)
	if err != nil {
		s.logger.Error("trip not found for update",
			slog.Uint64("trip_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	if req.FromCity != nil {
		trip.FromCity = *req.FromCity
	}
	if req.ToCity != nil {
		trip.ToCity = *req.ToCity
	}
	if req.StartTime != nil {
		trip.StartTime = *req.StartTime
	}
	if req.DurationMin != nil {
		trip.DurationMin = *req.DurationMin
	}
	if req.AvailableSeats != nil {
		trip.AvailableSeats = *req.AvailableSeats
	}
	if req.Price != nil {
		trip.Price = *req.Price
	}
	if req.TripStatus != nil {
		trip.TripStatus = string(*req.TripStatus)
	}

	if err := s.tripRepo.Update(trip); err != nil {
		s.logger.Error("failed to update trip",
			slog.Uint64("trip_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	return trip, nil
}

func (s *tripService) Delete(id uint) error {
	_, err := s.tripRepo.GetByID(id)
	if err != nil {
		s.logger.Error("trip not found for delete",
			slog.Uint64("trip_id", uint64(id)),
			slog.Any("error", err),
		)
		return err
	}

	if err := s.tripRepo.Delete(id); err != nil {
		s.logger.Error("failed to delete trip",
			slog.Uint64("trip_id", uint64(id)),
			slog.Any("error", err),
		)
		return err
	}

	return nil
}
