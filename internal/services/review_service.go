package services

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type ReviewService interface {
	Create(tripID uint, req *models.Review) (*models.Review, error)

	List() (*models.Review, error)

	GetByID(id uint) (*models.Review, error)

	Update(id uint, req *models.Review) (*models.Review, error)

	Delete(id uint) error
}
