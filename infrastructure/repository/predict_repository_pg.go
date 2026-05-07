// Esta struct implementa el PORT definido en domain/ports
// Go verifica esto en tiempo de compilación con la línea del var _
package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
	"time"
)

type predictRepositoryPG struct {
	db *sql.DB
}

// ✅ Go valida en compilación que implementa la interfaz
var _ ports.PredictionRepository = (*predictRepositoryPG)(nil)

func NewPredictionRepository(db *sql.DB) ports.PredictionRepository {
	return &predictRepositoryPG{db: db}
}

func (r *predictRepositoryPG) Save(prediction *models.Prediction) error {
	query := `
        INSERT INTO predictions (id, service_type, pieces_count, estimated_time, actual_time, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	var actualMins *float64
	if prediction.Actual != nil {
		mins := prediction.Actual.Minutes()
		actualMins = &mins
	}

	_, err := r.db.Exec(query,
		prediction.ID,
		prediction.ServiceType,
		prediction.PiecesCount,
		prediction.Estimated.Minutes(),
		actualMins,
		prediction.CreatedAt,
	)
	return err
}

func (r *predictRepositoryPG) FindByServiceType(serviceType string) ([]*models.Prediction, error) {
	query := `
		SELECT id, service_type, pieces_count, estimated_time, actual_time, created_at 
		FROM predictions 
		WHERE service_type = $1
	`
	rows, err := r.db.Query(query, serviceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []*models.Prediction
	for rows.Next() {
		var p models.Prediction
		var estMins float64
		var actMins sql.NullFloat64

		if err := rows.Scan(&p.ID, &p.ServiceType, &p.PiecesCount, &estMins, &actMins, &p.CreatedAt); err != nil {
			return nil, err
		}

		p.Estimated = time.Duration(estMins * float64(time.Minute))
		if actMins.Valid {
			actual := time.Duration(actMins.Float64 * float64(time.Minute))
			p.Actual = &actual
		}
		predictions = append(predictions, &p)
	}

	return predictions, rows.Err()
}

func (r *predictRepositoryPG) GetHistoricalData(serviceType string) ([]float64, error) {
	// Solo tomamos en cuenta los que tienen tiempo real (ya finalizados)
	query := `SELECT actual_time FROM predictions WHERE service_type = $1 AND actual_time IS NOT NULL`

	rows, err := r.db.Query(query, serviceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historical []float64
	for rows.Next() {
		var actMins float64
		if err := rows.Scan(&actMins); err != nil {
			return nil, err
		}
		historical = append(historical, actMins)
	}

	return historical, rows.Err()
}
