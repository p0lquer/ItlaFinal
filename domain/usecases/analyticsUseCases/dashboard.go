package analyticsUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
)

type DashboardUseCase struct{ repo ports.AnalyticsRepository }

func NewDashboardUseCase(repo ports.AnalyticsRepository) *DashboardUseCase {
	return &DashboardUseCase{repo: repo}
}
func (uc *DashboardUseCase) Execute() (*models.AnalyticsDashboard, error) { return uc.repo.Dashboard() }
