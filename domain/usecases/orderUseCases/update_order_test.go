package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/pkg/predictor"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type statusOrderRepo struct {
	order        *models.Order
	transition   bool
	transitionTo models.OrderStatus
}

func (r *statusOrderRepo) Create(*models.Order) error { return nil }
func (r *statusOrderRepo) FindByID(string) (*models.Order, error) {
	if r.order == nil {
		return nil, errors.New("not found")
	}
	return r.order, nil
}
func (r *statusOrderRepo) FindAll() ([]*models.Order, error)                { return nil, nil }
func (r *statusOrderRepo) FindByUserID(uuid.UUID) ([]*models.Order, error)  { return nil, nil }
func (r *statusOrderRepo) FindByCustomerID(string) ([]*models.Order, error) { return nil, nil }
func (r *statusOrderRepo) UpdateStatus(string, models.OrderStatus) error    { return nil }
func (r *statusOrderRepo) Delete(string) error                              { return nil }
func (r *statusOrderRepo) Transition(_ string, target models.OrderStatus) (bool, error) {
	r.transitionTo = target
	return r.transition, nil
}

type statusPredictionRepo struct{ saved *models.Prediction }

func (r *statusPredictionRepo) Save(*models.Prediction) error { return nil }
func (r *statusPredictionRepo) FindByServiceType(string) ([]*models.Prediction, error) {
	return nil, nil
}
func (r *statusPredictionRepo) GetHistoricalData(string) ([]predictor.DataPoint, error) {
	return nil, nil
}
func (r *statusPredictionRepo) UpdateActualTime(string, time.Duration) error { return nil }
func (r *statusPredictionRepo) UpsertForOrder(_ string, p *models.Prediction) error {
	r.saved = p
	return nil
}

type statusNotifier struct {
	statuses []string
	ready    int
}

func (n *statusNotifier) NotifyOrderReady(string, string) error { n.ready++; return nil }
func (n *statusNotifier) NotifyStatusChange(_ string, _ string, status string) error {
	n.statuses = append(n.statuses, status)
	return nil
}

var _ ports.OrderRepository = (*statusOrderRepo)(nil)
var _ ports.PredictionRepository = (*statusPredictionRepo)(nil)

func TestUpdateStatusReadyPersistsTrainingAndNotifies(t *testing.T) {
	started := time.Now().Add(-20 * time.Minute)
	orders := &statusOrderRepo{order: &models.Order{ID: "o-1", CustomerID: "c-1", ServiceType: "planchado", Status: models.StatusProcessing, Weight: 2, PiecesCount: 3, EstimatedTime: 10 * time.Minute, StartedAt: &started}, transition: true}
	predictions := &statusPredictionRepo{}
	notifier := &statusNotifier{}
	if err := NewUpdateOrderStatusUseCase(orders, predictions, notifier).Execute("o-1", models.StatusReady); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if predictions.saved == nil || predictions.saved.Actual == nil || *predictions.saved.Actual < 15*time.Minute {
		t.Fatalf("training prediction not saved correctly: %#v", predictions.saved)
	}
	if len(notifier.statuses) != 1 || notifier.statuses[0] != string(models.StatusReady) || notifier.ready != 1 {
		t.Fatalf("notifications = %#v ready=%d", notifier.statuses, notifier.ready)
	}
}

func TestUpdateStatusRejectsSkippedTransition(t *testing.T) {
	orders := &statusOrderRepo{order: &models.Order{ID: "o-1", Status: models.StatusReceived}, transition: true}
	err := NewUpdateOrderStatusUseCase(orders, &statusPredictionRepo{}).Execute("o-1", models.StatusReady)
	if err == nil || orders.transitionTo != "" {
		t.Fatalf("skipped transition was accepted: err=%v transition=%s", err, orders.transitionTo)
	}
}

func TestIsValidStatusTransition(t *testing.T) {
	if !IsValidStatusTransition(models.StatusReceived, models.StatusProcessing) || IsValidStatusTransition(models.StatusReceived, models.StatusReady) {
		t.Fatal("unexpected transition validation")
	}
}
