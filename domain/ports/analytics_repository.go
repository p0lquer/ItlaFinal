package ports

import "ITLAFINAL/domain/models"

type AnalyticsRepository interface {
	Dashboard() (*models.AnalyticsDashboard, error)
}
