package usecases

import (
	"ITLAFINAL/domain/ports"
)

type DeleteOrderUseCase struct {
	orderRepo ports.OrderRepository
	predRepo  ports.PredictionRepository
}

func NewDeleteOrderUseCase(
	orderRepo ports.OrderRepository,
	predRepo ports.PredictionRepository) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		orderRepo: orderRepo,
		predRepo:  predRepo,
	}
}

func (uc *DeleteOrderUseCase) Execute(orderID string) error {
	// 1. Obtener la orden para conocer su tipo de servicio
	err := uc.orderRepo.Delete(orderID)
	if err != nil {
		return err
	}
	return nil
}
