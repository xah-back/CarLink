package services

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/constants"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
)

type TripService interface {
	Create(driverID uint, req *models.TripCreateRequest) (*models.Trip, error)

	// List() ([]models.Trip, error)

	// GetByID(id uint) (*models.Trip, error)

	// Update(id uint, req *models.Trip) (*models.Trip, error)

	// Delete(id uint) error
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

func (s *tripService) Create(id uint, req *models.TripCreateRequest) (*models.Trip, error) {
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
		From:           req.From,
		To:             req.To,
		StartTime:      req.StartTime,
		DurationMin:    req.DurationMin,
		TotalSeats:     car.Seats,
		AvailableSeats: req.AvailableSeats,
		Price:          req.Price,
		TripStatus:     string(constants.TripPublished),
	}

	if err := s.tripRepo.Create(&trip); err != nil {
		return nil, err
	}

	return &trip, nil
}
