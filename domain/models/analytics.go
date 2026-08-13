package models

type AnalyticsPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type AnalyticsSummary struct {
	OrdersToday     int     `json:"orders_today"`
	InProcess       int     `json:"in_process"`
	ReadyForPickup  int     `json:"ready_for_pickup"`
	DeliveredToday  int     `json:"delivered_today"`
	ActiveCustomers int     `json:"active_customers"`
	PaidRevenue     float64 `json:"paid_revenue"`
}

type AnalyticsDashboard struct {
	Summary         AnalyticsSummary `json:"summary"`
	OrdersByStatus  []AnalyticsPoint `json:"orders_by_status"`
	OrdersByService []AnalyticsPoint `json:"orders_by_service"`
	OrdersTrend     []AnalyticsPoint `json:"orders_trend"`
	RevenueTrend    []AnalyticsPoint `json:"revenue_trend"`
	UsersByRole     []AnalyticsPoint `json:"users_by_role"`
}
