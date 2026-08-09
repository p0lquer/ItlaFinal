// Esta struct implementa el PORT definido en domain/ports
// Go verifica esto en tiempo de compilación con la línea del var _
package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/pkg/predictor"
	"database/sql"
	"fmt"
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
        INSERT INTO predictions (id, order_id, service_type, pieces_count, estimated, actual, weight, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	var actualMins *float64
	if prediction.Actual != nil {
		mins := prediction.Actual.Minutes()
		actualMins = &mins
	}

	_, err := r.db.Exec(query,
		prediction.ID,
		prediction.OrderID,
		prediction.ServiceType,
		prediction.PiecesCount,
		prediction.Estimated.Minutes(),
		actualMins,
		prediction.Weight,
		prediction.CreatedAt,
	)
	return err
}

func (r *predictRepositoryPG) FindByServiceType(serviceType string) ([]*models.Prediction, error) {
	query := `
		SELECT id, service_type, pieces_count, estimated_time, actual, created_at 
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

func (r *predictRepositoryPG) GetHistoricalData(serviceType string) ([]predictor.DataPoint, error) {
	// Solo tomamos en cuenta los que tienen tiempo real (ya finalizados) Y peso
	// registrado — órdenes viejas, previas a esta migración, tendrán weight NULL
	// y no sirven como punto de entrenamiento para la regresión por peso.
	query := `
		SELECT weight, actual
		FROM predictions
		WHERE service_type = $1 AND actual IS NOT NULL AND weight IS NOT NULL
	`
	rows, err := r.db.Query(query, serviceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historical []predictor.DataPoint
	for rows.Next() {
		var weight, actMins float64
		if err := rows.Scan(&weight, &actMins); err != nil {
			return nil, err
		}
		historical = append(historical, predictor.DataPoint{Weight: weight, ActualMinutes: actMins})
	}
	return historical, rows.Err()
}

func (r *predictRepositoryPG) UpdateActualTime(
	orderID string,
	actual time.Duration,
) error {

	actualMinutes := actual.Minutes()

	_, err := r.db.Exec(`
        UPDATE predictions
        SET actual = $1
        WHERE order_id = $2
    `,
		actualMinutes,
		orderID,
	)

	return err
}

// UpsertForOrder escribe (o actualiza) la única fila de entrenamiento por
// orden. ON CONFLICT (order_id) impide duplicar datos al re-ejecutar la
// transición a "lista" (carreras worker/manual incluidas).
func (r *predictRepositoryPG) UpsertForOrder(orderID string, p *models.Prediction) error {
	_, err := r.db.Exec(`
		INSERT INTO predictions (id, order_id, service_type, pieces_count,
		                         estimated_time, actual_time, weight, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_id) DO UPDATE SET
			estimated_time = EXCLUDED.estimated_time,
			actual_time    = EXCLUDED.actual_time,
			pieces_count   = EXCLUDED.pieces_count,
			weight         = EXCLUDED.weight`,
		p.ID, orderID, p.ServiceType, p.PiecesCount,
		p.Estimated.Nanoseconds(), durationToNullableNs(p.Actual), p.Weight, p.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert prediction for order: %w", err)
	}
	return nil
}
func durationToNullableNs(d *time.Duration) any {
	if d == nil {
		return nil
	}
	return d.Nanoseconds()
}
