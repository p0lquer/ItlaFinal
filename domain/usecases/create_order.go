package usecases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"time"

	"github.com/google/uuid"
)

type CreateOrderUseCase struct {
	orderRepo ports.OrderRepository
	predRepo  ports.PredictionRepository
}

func NewCreateOrderUseCase(
	orderRepo ports.OrderRepository,
	predRepo ports.PredictionRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo: orderRepo,
		predRepo:  predRepo,
	}
}

func (uc *CreateOrderUseCase) Execute(
	customerID, serviceType string,
	piecesCount int,
	notes string,
) (*models.Order, error) {

	// 1. Obtener datos históricos para predecir
	historicalData, err := uc.predRepo.GetHistoricalData(serviceType)
	if err != nil || len(historicalData) == 0 {
		// Si no hay historial, usar estimado por defecto según tipo
		historicalData = defaultEstimates(serviceType)
	}

	// 2. Calcular predicción
	estimatedMinutes := weightedAverage(historicalData)
	estimated := time.Duration(estimatedMinutes) * time.Minute

	// 3. Crear la orden
	order := &models.Order{
		ID:            uuid.NewString(),
		CustomerID:    customerID,
		ServiceType:   serviceType,
		PiecesCount:   piecesCount,
		Notes:         notes,
		Status:        models.StatusReceived,
		EstimatedTime: estimated,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 4. Persistir
	if err := uc.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func defaultEstimates(serviceType string) []float64 {
	defaults := map[string][]float64{
		"lavado_secado":  {60, 60, 60},
		"planchado":      {30, 30, 30},
		"lavado_en_seco": {120, 120, 120},
	}
	if v, ok := defaults[serviceType]; ok {
		return v
	}
	return []float64{60}
}

func weightedAverage(values []float64) float64 {
	if len(values) == 0 {
		return 60
	}

	var total, weightSum float64
	for i, v := range values {
		w := float64(i + 1)
		total += v * w
		weightSum += w
	}

	if weightSum == 0 {
		return 60
	}

	return total / weightSum
}
