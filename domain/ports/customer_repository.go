package ports

import "ITLAFINAL/domain/models"

type CustomerRepository interface {
	Create(customer *models.Customer) error
	FindAll() ([]*models.Customer, error)
	UpdateStatus(id string, status models.OrderStatus) error
	Delete(id string) error
}
