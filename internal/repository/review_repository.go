package repository

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type ReviewRepository interface {
	Create(review *models.Review) error

	List() ([]models.Review, error)

	GetByID(id uint) (*models.Review, error)

	Update(review *models.Review) (*models.Review, error)

	Delete(id uint) error
}
