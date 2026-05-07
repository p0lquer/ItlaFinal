package ports

import "ITLAFINAL/domain/models"

type PredictionRepository interface {
	Save(p *models.Prediction) error
	FindByServiceType(serviceType string) ([]*models.Prediction, error)
	GetHistoricalData(serviceType string) ([]float64, error) // retorna tiempos reales en minutos
}
