package ports

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/pkg/predictor"
	"time"
)

type PredictionRepository interface {
	Save(p *models.Prediction) error
	FindByServiceType(serviceType string) ([]*models.Prediction, error)
	// GetHistoricalData retorna los puntos (peso, tiempo real) YA filtrados
	// por tipo de servicio, listos para pkg/predictor.LinearRegression.
	GetHistoricalData(serviceType string) ([]predictor.DataPoint, error) // retorna tiempos reales en minutos
	UpdateActualTime(orderID string, actual time.Duration) error
}
