package customerUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/pkg/inputvalidation"

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

func (uc *CreateCustomerUseCase) Execute(_ string, name, phone, email string) (*models.Customer, error) {
	var err error
	name, err = inputvalidation.Name(name)
	if err != nil {
		return nil, err
	}
	phone, err = inputvalidation.Phone(phone)
	if err != nil {
		return nil, err
	}
	email, err = inputvalidation.Email(email)
	if err != nil {
		return nil, err
	}
	customer := &models.Customer{
		ID:    uuid.NewString(),
		Name:  name,
		Phone: &phone,
		Email: &email,
	}
	if err := uc.customerRepo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}
