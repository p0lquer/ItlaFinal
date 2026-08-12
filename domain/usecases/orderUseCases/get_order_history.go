package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type GetOrderHistoryUseCase struct {
	orderRepo ports.OrderRepository
}

func NewGetOrderHistoryUseCase(orderRepo ports.OrderRepository) *GetOrderHistoryUseCase {
	return &GetOrderHistoryUseCase{orderRepo: orderRepo}
}

func (uc *GetOrderHistoryUseCase) Execute(orderID string) ([]*models.OrderStatusChange, error) {
	return uc.orderRepo.FindStatusHistory(orderID)
}
