package customerUseCases

import (
	"ITLAFINAL/domain/ports"
	"errors"
)

type DeleteCustomerUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewDeleteCustomerUseCase(customerRepo ports.CustomerRepository) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{
		customerRepo: customerRepo,
	}
}

func (uc *DeleteCustomerUseCase) Execute(customerID string) error {
	// 1. Verificar que el cliente exista
	customer, err := uc.customerRepo.FindByID(customerID)
	if err != nil {
		return errors.New("cliente no encontrado")
	}

	// 2. Eliminar el cliente
	if err := uc.customerRepo.Delete(customer.ID); err != nil {
		return err
	}

	return nil
}
