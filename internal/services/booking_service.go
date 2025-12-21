package services

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type BookingService interface {
	Create(id uint, req *models.Booking) (*models.Booking, error)

	List() ([]models.Booking, error)

	GetByID(id uint) (*models.Booking, error)

	Update(id uint, req *models.Booking) (*models.Booking, error)

	Delete(id uint) error
}
