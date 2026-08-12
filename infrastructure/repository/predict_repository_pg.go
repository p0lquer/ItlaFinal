package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/pkg/predictor"
	"database/sql"
	"fmt"
	"time"
)

type predictRepositoryPG struct{ db *sql.DB }

var _ ports.PredictionRepository = (*predictRepositoryPG)(nil)

func NewPredictionRepository(db *sql.DB) ports.PredictionRepository {
	return &predictRepositoryPG{db: db}
}

// SQL durations use minutes, never Duration.Nanoseconds().
func (r *predictRepositoryPG) Save(p *models.Prediction) error {
	_, err := r.db.Exec(`INSERT INTO predictions
		(id, order_id, service_type, pieces_count, estimated_time, actual_time, weight, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, p.ID, p.OrderID, p.ServiceType,
		p.PiecesCount, p.Estimated.Minutes(), nullableMinutes(p.Actual), p.Weight, p.CreatedAt)
	return err
}

func (r *predictRepositoryPG) FindByServiceType(serviceType string) ([]*models.Prediction, error) {
	rows, err := r.db.Query(`SELECT id, order_id, service_type, pieces_count, estimated_time, actual_time, weight, created_at
		FROM predictions WHERE service_type=$1 ORDER BY created_at DESC`, serviceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*models.Prediction, 0)
	for rows.Next() {
		var p models.Prediction
		var estimated float64
		var actual sql.NullFloat64
		var weight sql.NullFloat64
		if err := rows.Scan(&p.ID, &p.OrderID, &p.ServiceType, &p.PiecesCount, &estimated, &actual, &weight, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.Estimated = time.Duration(estimated * float64(time.Minute))
		if actual.Valid {
			d := time.Duration(actual.Float64 * float64(time.Minute))
			p.Actual = &d
		}
		if weight.Valid {
			p.Weight = weight.Float64
		}
		result = append(result, &p)
	}
	return result, rows.Err()
}

func (r *predictRepositoryPG) GetHistoricalData(serviceType string) ([]predictor.DataPoint, error) {
	rows, err := r.db.Query(`SELECT weight, actual_time FROM predictions
		WHERE service_type=$1 AND actual_time IS NOT NULL AND weight IS NOT NULL`, serviceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := make([]predictor.DataPoint, 0)
	for rows.Next() {
		var point predictor.DataPoint
		if err := rows.Scan(&point.Weight, &point.ActualMinutes); err != nil {
			return nil, err
		}
		data = append(data, point)
	}
	return data, rows.Err()
}

func (r *predictRepositoryPG) UpdateActualTime(orderID string, actual time.Duration) error {
	_, err := r.db.Exec(`UPDATE predictions SET actual_time=$1 WHERE order_id=$2`, actual.Minutes(), orderID)
	return err
}

func (r *predictRepositoryPG) UpsertForOrder(orderID string, p *models.Prediction) error {
	_, err := r.db.Exec(`INSERT INTO predictions (id,order_id,service_type,pieces_count,estimated_time,actual_time,weight,created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (order_id) DO UPDATE SET service_type=EXCLUDED.service_type,
		pieces_count=EXCLUDED.pieces_count, estimated_time=EXCLUDED.estimated_time,
		actual_time=EXCLUDED.actual_time, weight=EXCLUDED.weight`,
		p.ID, orderID, p.ServiceType, p.PiecesCount, p.Estimated.Minutes(), nullableMinutes(p.Actual), p.Weight, p.CreatedAt)
	if err != nil {
		return fmt.Errorf("upsert prediction for order: %w", err)
	}
	return nil
}

func nullableMinutes(d *time.Duration) any {
	if d == nil {
		return nil
	}
	return d.Minutes()
}
