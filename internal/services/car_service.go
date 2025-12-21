package services

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type CarService interface {
	Create(id uint, req *models.Car) (*models.Car, error)

	List() ([]models.Car, error)

	GetByID(id uint) (*models.Car, error)

	Update(id uint, req *models.Car) (*models.Car, error)

	Delete(id uint) error
}
