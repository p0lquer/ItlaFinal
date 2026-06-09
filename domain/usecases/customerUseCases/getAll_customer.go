package customerUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type GetAllCustomersUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewGetAllCustomersUseCase(customerRepo ports.CustomerRepository) *GetAllCustomersUseCase {
	return &GetAllCustomersUseCase{customerRepo: customerRepo}
}

func (uc *GetAllCustomersUseCase) Execute() ([]*models.Customer, error) {
	return uc.customerRepo.FindAll()
}

func (uc *GetAllCustomersUseCase) ExecuteByID(customerID string) (*models.Customer, error) {
	return uc.customerRepo.FindByID(customerID)
}
