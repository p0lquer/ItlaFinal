package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/pkg/predictor"
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
	weight float64,
) (*models.Order, error) {

	// 1. Obtener datos históricos para predecir
	historicalData, err := uc.predRepo.GetHistoricalData(serviceType)
	if err != nil {
		// Si no hay historial, usar estimado por defecto según tipo
		return nil, err
	}

	// 2. Calcular predicción
	var estimatedMinutes float64
	if len(historicalData) >= 2 {
		predict := predictor.LinearRegression(historicalData)
		estimatedMinutes = predict(weight)
	} else {
		estimatedMinutes = defaultEstimate(serviceType)
	}
	estimated := time.Duration(estimatedMinutes) * time.Minute

	order := &models.Order{
		ID: uuid.NewString(), CustomerID: customerID, ServiceType: serviceType,
		PiecesCount: piecesCount, Weight: weight, Notes: notes,
		Status: models.StatusReceived, EstimatedTime: estimated,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if err := uc.orderRepo.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

func defaultEstimate(serviceType string) float64 {
	defaults := map[string]float64{"lavado_secado": 60, "planchado": 30, "lavado_en_seco": 120}
	if v, ok := defaults[serviceType]; ok {
		return v
	}
	return 60
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
