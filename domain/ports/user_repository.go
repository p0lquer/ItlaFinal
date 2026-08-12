package ports

import "ITLAFINAL/domain/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	Delete(id string) error
	List(filter models.UserFilter) (*models.UserPage, error)
	SetActive(id string, active bool) error
	DeleteManaged(id string) error
	EnsureBootstrapAdmin(user *models.User) error
}
