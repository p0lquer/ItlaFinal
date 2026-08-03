package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type GetMyOrdersUseCase struct {
	orderRepo ports.OrderRepository
}

func NewGetMyOrdersUseCase(orderRepo ports.OrderRepository) *GetMyOrdersUseCase {
	return &GetMyOrdersUseCase{orderRepo: orderRepo}
}

func (uc *GetMyOrdersUseCase) Execute(customerID string) ([]*models.Order, error) {
	if customerID == "" {
		return nil, nil
	}
	return uc.orderRepo.FindByCustomerID(customerID)
}
