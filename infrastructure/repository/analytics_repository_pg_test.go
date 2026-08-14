package repository

import (
	"strings"
	"testing"
)

func TestAnalyticsTrendQueriesAreBoundedAndPreserveEmptyDays(t *testing.T) {
	tests := []struct {
		name  string
		query string
		field string
	}{
		{name: "orders", query: ordersTrendQuery, field: "created_at"},
		{name: "revenue", query: revenueTrendQuery, field: "paid_at"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !strings.Contains(test.query, "generate_series(CURRENT_DATE - INTERVAL '6 day', CURRENT_DATE, '1 day')") {
				t.Fatal("the seven calendar days must remain present, including zero-value days")
			}
			if !strings.Contains(test.query, test.field+" >= CURRENT_DATE - INTERVAL '6 day'") || !strings.Contains(test.query, test.field+" < CURRENT_DATE + INTERVAL '1 day'") {
				t.Fatalf("expected a sargable, bounded %s range", test.field)
			}
		})
	}
}
