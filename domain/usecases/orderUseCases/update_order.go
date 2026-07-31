package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"time"

	"github.com/google/uuid"
)

type UpdateOrderStatusUseCase struct {
	orderRepo ports.OrderRepository
	predRepo  ports.PredictionRepository
}

func NewUpdateOrderStatusUseCase(orderRepo ports.OrderRepository, predRepo ports.PredictionRepository) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{orderRepo: orderRepo, predRepo: predRepo}
}

func (uc *UpdateOrderStatusUseCase) Execute(orderID string, status models.OrderStatus) error {
	if err := uc.orderRepo.UpdateStatus(orderID, status); err != nil {
		return err
	}

	// Cuando la orden pasa a "lista", ya conocemos el tiempo REAL que tomó.
	// Se guarda como dato de entrenamiento para la próxima predicción de este servicio.
	if status == models.StatusReady {
		order, err := uc.orderRepo.FindByID(orderID)
		if err != nil {
			return err
		}
		actual := time.Since(order.CreatedAt)
		_ = uc.predRepo.Save(&models.Prediction{
			ID: uuid.NewString(), ServiceType: order.ServiceType, PiecesCount: order.PiecesCount,
			Weight: order.Weight, Estimated: order.EstimatedTime, Actual: &actual, CreatedAt: time.Now(),
		})
	}
	return nil
}
