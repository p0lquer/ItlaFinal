package customerUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"time"

	"github.com/google/uuid"
)

type CreateCustomerUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewCreateCustomerUseCase(customerRepo ports.CustomerRepository) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		customerRepo: customerRepo,
	}
}

func (uc *CreateCustomerUseCase) Execute(name, phone, email string) (*models.Customer, error) {
	customer := &models.Customer{
		ID:        uuid.NewString(),
		Name:      name,
		Phone:     phone,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := uc.customerRepo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}
