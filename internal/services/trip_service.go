package services

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type TripService interface {
	Create(driverID uint, req *models.Trip) (*models.Trip, error)

	List() ([]models.Trip, error)

	GetByID(id uint) (*models.Trip, error)

	Update(id uint, req *models.Trip) (*models.Trip, error)

	Delete(id uint) error
}
