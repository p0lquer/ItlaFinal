package ports

import "ITLAFINAL/domain/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	Delete(id string) error
}
