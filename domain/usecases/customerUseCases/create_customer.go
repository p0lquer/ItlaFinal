package customerUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type CreateCustomerUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewCreateCustomerUseCase(customerRepo ports.CustomerRepository) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		customerRepo: customerRepo,
	}
}

func (uc *CreateCustomerUseCase) Execute(id, name, phone, email string) (*models.Customer, error) {
	customer := &models.Customer{
		ID:    id,
		Name:  name,
		Phone: &phone,
		Email: &email,
	}
	if err := uc.customerRepo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}
