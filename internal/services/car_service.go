package services

import (
	"log/slog"

	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"github.com/mutsaevz/team-5-ambitious/internal/repository"
)

type CarService interface {
	Create(id uint, req models.CarCreateRequest) (*models.Car, error)

	// List() ([]models.Car, error)

	GetByID(id uint) (*models.Car, error)

	// Update(id uint, req *models.Car) (*models.Car, error)

	// Delete(id uint) error
}

type carService struct {
	carRepo  repository.CarRepository
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewCarService(carRepo repository.CarRepository, userRepo repository.UserRepository, logger *slog.Logger) CarService {
	return &carService{
		carRepo:  carRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *carService) Create(id uint, req models.CarCreateRequest) (*models.Car, error) {
	driver, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	var car = models.Car{
		OwnerID:  driver.ID,
		Brand:    req.Brand,
		CarModel: req.CarModel,
		Seats:    req.Seats,
	}

	if err := s.carRepo.Create(&car); err != nil {
		return nil, err
	}

	return &car, err
}

func (s *carService) GetByID(id uint) (*models.Car, error) {
	car, err := s.carRepo.GetByID(id)

	if err != nil {
		return nil, err
	}

	return car, nil
}
