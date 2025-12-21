package services

import "github.com/mutsaevz/team-5-ambitious/internal/models"

type UserService interface {
	Create(req *models.User) (*models.User, error)

	List() ([]models.User, error)

	GetByID(id uint) (*models.User, error)

	Update(id uint, req models.User) (*models.User, error)

	Delete(id uint) error
}
