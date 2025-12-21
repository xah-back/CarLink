package repository

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type BookingRepository interface {
	Create(booking *models.Booking) error

	List() ([]models.Booking, error)

	GetByID(id uint) (*models.Booking, error)

	Update(booking *models.Booking) (*models.Booking, error)

	Delete(id uint) error
}
