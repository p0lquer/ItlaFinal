package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/pkg/predictor"
	"log"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateOrderUseCase struct {
	orderRepo       ports.OrderRepository
	predRepo        ports.PredictionRepository
	serviceTypeRepo ports.ServiceTypeRepository
}

func NewCreateOrderUseCase(
	orderRepo ports.OrderRepository,
	predRepo ports.PredictionRepository,
	serviceTypeRepo ports.ServiceTypeRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:       orderRepo,
		predRepo:        predRepo,
		serviceTypeRepo: serviceTypeRepo,
	}
}

func (uc *CreateOrderUseCase) Execute(
	customerID, serviceType string,
	piecesCount int,
	notes string,
	weight float64,
) (*models.Order, error) {
	serviceKey := normalizeServiceType(serviceType)
	serviceConfig, _ := uc.serviceTypeRepo.FindByName(serviceType)
	if serviceConfig == nil {
		serviceConfig, _ = uc.serviceTypeRepo.FindByName(serviceKey)
	}

	// // 1. Obtener datos históricos para predecir
	// historicalData, err := uc.predRepo.GetHistoricalData(serviceKey)
	// if err != nil {
	// 	// Si no hay historial, usar estimado por defecto según tipo
	// 	return nil, err
	// }
	// if len(historicalData) == 0 {
	// Si no hay historial, usar estimado por defecto según tipo
	estimatedMinutes := uc.predictMinutes(serviceKey, weight)
	estimated := time.Duration(estimatedMinutes) * time.Minute
	estimatedCost := estimateCost(serviceConfig, piecesCount, weight, serviceKey)

	order := &models.Order{
		ID: uuid.NewString(), CustomerID: customerID, ServiceType: serviceType,
		PiecesCount: piecesCount, Weight: weight, Notes: notes,
		Status: models.StatusReceived, EstimatedTime: estimated, EstimatedCost: estimatedCost,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	if err := uc.orderRepo.Create(order); err != nil {
		return nil, err
	}
	// 2. Ya NO se guarda Prediction aquí. La única fila de entrenamiento
	//    se crea al pasar a "lista" (ver UpdateOrderStatusUseCase).
	return order, nil
}

// predictMinutes usa regresión si hay >=2 puntos del MISMO tipo; si no,
// cae al estimado por defecto. Un error de repo solo se registra.
func (uc *CreateOrderUseCase) predictMinutes(serviceKey string, weight float64) float64 {
	historicalData, err := uc.predRepo.GetHistoricalData(serviceKey)
	if err != nil {
		log.Printf("create_order: GetHistoricalData(%q) error=%v; usando default", serviceKey, err)
		return defaultEstimate(serviceKey)
	}
	if len(historicalData) >= 2 {
		return predictor.LinearRegression(historicalData)(weight)
	}
	return defaultEstimate(serviceKey)
}

func defaultEstimate(serviceType string) float64 {
	type serviceDefaults struct {
		minutes                   float64
		base, perWeight, perPiece float64
	}

	defaults := map[string]serviceDefaults{
		"lavado_y_secado": {minutes: 60, base: 50, perWeight: 5, perPiece: 0},
		"planchado":       {minutes: 30, base: 30, perWeight: 0, perPiece: 2},
		"lavado_en_seco":  {minutes: 120, base: 80, perWeight: 8, perPiece: 0},
	}
	if v, ok := defaults[serviceType]; ok {
		return v.minutes
	}
	return 60
}

func normalizeServiceType(serviceType string) string {
	serviceType = strings.TrimSpace(strings.ToLower(serviceType))
	serviceType = strings.NewReplacer(
		"á", "a",
		"é", "e",
		"í", "i",
		"ó", "o",
		"ú", "u",
		" ", "_",
		"-", "_",
	).Replace(serviceType)
	return serviceType
}

func estimateCost(serviceType *models.ServiceType, piecesCount int, weight float64, serviceKey string) float64 {
	if serviceType != nil {
		cost := serviceType.BasePrice
		cost += serviceType.PricePerWeight * weight
		cost += serviceType.PricePerPiece * float64(piecesCount)
		if cost > 0 {
			return roundMoney(cost)
		}
	}

	defaults := map[string]struct {
		base      float64
		perWeight float64
		perPiece  float64
	}{
		"lavado_secado":   {base: 50, perWeight: 5},
		"lavado_y_secado": {base: 50, perWeight: 5},
		"planchado":       {base: 30, perPiece: 2},
		"lavado_en_seco":  {base: 80, perWeight: 8},
	}
	if v, ok := defaults[serviceKey]; ok {
		return roundMoney(v.base + (v.perWeight * weight) + (v.perPiece * float64(piecesCount)))
	}
	return roundMoney((5 * weight) + float64(piecesCount))
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
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
