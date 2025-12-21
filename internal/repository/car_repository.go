package repository

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type CarRepository interface {
	Create(car *models.Car) error

	List() ([]models.Car, error)

	GetByID(id uint) (*models.Car, error)

	Update(car *models.Car) (*models.Car, error)

	Delete(id uint) error
}
