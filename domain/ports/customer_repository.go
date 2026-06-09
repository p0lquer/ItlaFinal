package ports

import "ITLAFINAL/domain/models"

type CustomerRepository interface {
	Create(customer *models.Customer) error
	FindAll() ([]*models.Customer, error)
	FindByID(customerID string) (*models.Customer, error)
}
