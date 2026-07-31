package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"

	"github.com/google/uuid"
)

type GetMineOrdersUseCase struct {
	orderRepo ports.OrderRepository
	userID    uuid.UUID
}

func NewGetMineOrdersUseCase(orderRepo ports.OrderRepository, userID uuid.UUID) *GetMineOrdersUseCase {
	return &GetMineOrdersUseCase{orderRepo: orderRepo, userID: userID}
}

func (uc *GetMineOrdersUseCase) Execute(userID uuid.UUID) ([]*models.Order, error) {
	return uc.orderRepo.FindByUserID(uc.userID)
}
