package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimitRejectsExcessiveRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/login", RateLimit(2, time.Minute), func(c *gin.Context) { c.Status(http.StatusNoContent) })
	for attempt := 1; attempt <= 3; attempt++ {
		request := httptest.NewRequest(http.MethodGet, "/login", nil)
		request.RemoteAddr = "198.51.100.10:1234"
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if attempt < 3 && response.Code != http.StatusNoContent {
			t.Fatalf("attempt %d status=%d, want %d", attempt, response.Code, http.StatusNoContent)
		}
		if attempt == 3 && response.Code != http.StatusTooManyRequests {
			t.Fatalf("attempt %d status=%d, want %d", attempt, response.Code, http.StatusTooManyRequests)
		}
	}
}
