package repository

import (
	"testing"
	"time"
)

func TestNullableMinutes(t *testing.T) {
	if got := nullableMinutes(nil); got != nil {
		t.Fatalf("nil duration = %#v, want nil", got)
	}
	duration := 90*time.Minute + 30*time.Second
	got, ok := nullableMinutes(&duration).(float64)
	if !ok || got != 90.5 {
		t.Fatalf("duration stored as %#v, want 90.5 minutes", nullableMinutes(&duration))
	}
}
