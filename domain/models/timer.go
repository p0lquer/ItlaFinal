package models

import "time"

type TimerStatus string

const (
	TimerRunning  TimerStatus = "corriendo"
	TimerPaused   TimerStatus = "pausado"
	TimerFinished TimerStatus = "finalizado"
)

type Timer struct {
	ID        string
	OrderID   string
	StartedAt time.Time
	EndsAt    time.Time
	Status    TimerStatus
}
