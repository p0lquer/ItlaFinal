package handlers

import (
	"ITLAFINAL/domain/usecases/analyticsUseCases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AnalyticsHandler struct {
	dashboard *analyticsUseCases.DashboardUseCase
}

func NewAnalyticsHandler(dashboard *analyticsUseCases.DashboardUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{dashboard: dashboard}
}
func (h *AnalyticsHandler) Dashboard(c *gin.Context) {
	data, err := h.dashboard.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron cargar las estadísticas"})
		return
	}
	c.JSON(http.StatusOK, data)
}
