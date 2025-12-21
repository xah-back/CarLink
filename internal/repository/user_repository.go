package repository

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type UserRepository interface {
	Create(user *models.User) error

	List() ([]models.User, error)

	GetByID(id uint) (*models.User, error)

	Update(user *models.User) (*models.User, error)

	Delete(id uint) error
}
