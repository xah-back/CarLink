package repository

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type TripRepository interface {
	Create(trip *models.Trip) error

	List() ([]models.Trip, error)

	GetByID(id uint) (*models.Trip, error)

	Update(trip *models.Trip) (*models.Trip, error)

	Delete(id uint) error
}
