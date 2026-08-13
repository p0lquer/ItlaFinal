package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
)

type analyticsRepositoryPG struct{ db *sql.DB }

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
	if dashboard.OrdersTrend, err = r.points(`SELECT TO_CHAR(day, 'DD Mon'), COALESCE(COUNT(o.id),0)::float8 FROM generate_series(CURRENT_DATE - INTERVAL '6 day', CURRENT_DATE, '1 day') day LEFT JOIN orders o ON o.created_at >= day AND o.created_at < day + INTERVAL '1 day' GROUP BY day ORDER BY day`); err != nil {
		return nil, err
	}
	if dashboard.RevenueTrend, err = r.points(`SELECT TO_CHAR(day, 'DD Mon'), COALESCE(SUM(p.amount),0)::float8 FROM generate_series(CURRENT_DATE - INTERVAL '6 day', CURRENT_DATE, '1 day') day LEFT JOIN payments p ON p.paid_at >= day AND p.paid_at < day + INTERVAL '1 day' AND p.status='paid' GROUP BY day ORDER BY day`); err != nil {
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
