package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"time"

	"github.com/google/uuid"
)

const minRecordedMinutes = 15 * time.Minute

type UpdateOrderStatusUseCase struct {
	orderRepo ports.OrderRepository
	predRepo  ports.PredictionRepository
}

func NewUpdateOrderStatusUseCase(orderRepo ports.OrderRepository, predRepo ports.PredictionRepository) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{orderRepo: orderRepo, predRepo: predRepo}
}

func (uc *UpdateOrderStatusUseCase) Execute(orderID string, target models.OrderStatus) error {
	updated, err := uc.orderRepo.Transition(orderID, target)
	if err != nil {
		return err
	}
	if !updated || target != models.StatusReady {
		return nil
	}
	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	actual := order.ElapsedTime()
	if actual < minRecordedMinutes {
		actual = minRecordedMinutes
	}

	prediction := &models.Prediction{
		ID:          uuid.NewString(),
		OrderID:     order.ID,
		ServiceType: normalizeServiceType(order.ServiceType), // misma clave que GetHistoricalData
		PiecesCount: order.PiecesCount,
		Weight:      order.Weight,
		Estimated:   order.EstimatedTime,
		Actual:      &actual,
		CreatedAt:   time.Now(),
	}

	return uc.predRepo.UpsertForOrder(order.ID, prediction)
}

// 		if err := uc.predRepo.UpdateActualTime(orderID, actual); err != nil {
// 			return err
// 		}

// 		_ = uc.predRepo.Save(&models.Prediction{
// 			ID: uuid.NewString(), ServiceType: order.ServiceType, PiecesCount: order.PiecesCount,
// 			Weight: order.Weight, Estimated: order.EstimatedTime, Actual: &actual, CreatedAt: time.Now(),
// 		})
// 	}
// 	return nil
// }
