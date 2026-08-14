package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
)

type analyticsRepositoryPG struct{ db *sql.DB }

const ordersTrendQuery = `WITH daily AS (
	SELECT created_at::date AS day, COUNT(*)::float8 AS value
	FROM orders
	WHERE created_at >= CURRENT_DATE - INTERVAL '6 day'
	  AND created_at < CURRENT_DATE + INTERVAL '1 day'
	GROUP BY created_at::date
)
SELECT TO_CHAR(days.day, 'DD Mon'), COALESCE(daily.value, 0)
FROM generate_series(CURRENT_DATE - INTERVAL '6 day', CURRENT_DATE, '1 day') AS days(day)
LEFT JOIN daily ON daily.day = days.day::date
ORDER BY days.day`

const revenueTrendQuery = `WITH daily AS (
	SELECT paid_at::date AS day, SUM(amount)::float8 AS value
	FROM payments
	WHERE status = 'paid'
	  AND paid_at >= CURRENT_DATE - INTERVAL '6 day'
	  AND paid_at < CURRENT_DATE + INTERVAL '1 day'
	GROUP BY paid_at::date
)
SELECT TO_CHAR(days.day, 'DD Mon'), COALESCE(daily.value, 0)
FROM generate_series(CURRENT_DATE - INTERVAL '6 day', CURRENT_DATE, '1 day') AS days(day)
LEFT JOIN daily ON daily.day = days.day::date
ORDER BY days.day`

var _ ports.AnalyticsRepository = (*analyticsRepositoryPG)(nil)

func NewAnalyticsRepository(db *sql.DB) ports.AnalyticsRepository {
	return &analyticsRepositoryPG{db: db}
}

func (r *analyticsRepositoryPG) Dashboard() (*models.AnalyticsDashboard, error) {
	dashboard := &models.AnalyticsDashboard{}
	if err := r.db.QueryRow(`SELECT
		COUNT(*) FILTER (WHERE created_at >= CURRENT_DATE),
		COUNT(*) FILTER (WHERE status = 'en_proceso'),
		COUNT(*) FILTER (WHERE status = 'lista'),
		COUNT(*) FILTER (WHERE status = 'entregada' AND updated_at >= CURRENT_DATE),
		(SELECT COUNT(*) FROM users WHERE role='customer' AND is_active),
		COALESCE((SELECT SUM(amount) FROM payments WHERE status='paid'), 0)
		FROM orders`).Scan(&dashboard.Summary.OrdersToday, &dashboard.Summary.InProcess, &dashboard.Summary.ReadyForPickup, &dashboard.Summary.DeliveredToday, &dashboard.Summary.ActiveCustomers, &dashboard.Summary.PaidRevenue); err != nil {
		return nil, err
	}
	var err error
	if dashboard.OrdersByStatus, err = r.points(`SELECT status, COUNT(*)::float8 FROM orders GROUP BY status ORDER BY status`); err != nil {
		return nil, err
	}
	if dashboard.OrdersByService, err = r.points(`SELECT service_type, COUNT(*)::float8 FROM orders GROUP BY service_type ORDER BY COUNT(*) DESC, service_type LIMIT 6`); err != nil {
		return nil, err
	}
	if dashboard.UsersByRole, err = r.points(`SELECT role, COUNT(*)::float8 FROM users GROUP BY role ORDER BY role`); err != nil {
		return nil, err
	}
	// Aggregate the bounded seven-day range before joining its calendar. Besides
	// keeping zero-order days visible, the sargable created_at predicate lets
	// PostgreSQL use orders_created_at_idx instead of scanning all history.
	if dashboard.OrdersTrend, err = r.points(ordersTrendQuery); err != nil {
		return nil, err
	}
	// This is equivalent to the former calendar LEFT JOIN but reads only paid
	// payments in the displayed window; the partial covering index supports it.
	if dashboard.RevenueTrend, err = r.points(revenueTrendQuery); err != nil {
		return nil, err
	}
	return dashboard, nil
}

func (r *analyticsRepositoryPG) points(query string) ([]models.AnalyticsPoint, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	points := make([]models.AnalyticsPoint, 0)
	for rows.Next() {
		var point models.AnalyticsPoint
		if err := rows.Scan(&point.Label, &point.Value); err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, rows.Err()
}
