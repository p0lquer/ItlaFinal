package workers

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/domain/usecases/orderUseCases"
	"ITLAFINAL/pkg/predictor"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeOrderRepo struct {
	orders []*models.Order
}

func (f *fakeOrderRepo) Create(order *models.Order) error { return nil }
func (f *fakeOrderRepo) FindByID(id string) (*models.Order, error) { return nil, nil }
func (f *fakeOrderRepo) FindAll() ([]*models.Order, error) { return f.orders, nil }
func (f *fakeOrderRepo) FindByUserID(userID uuid.UUID) ([]*models.Order, error) { return nil, nil }
func (f *fakeOrderRepo) FindByCustomerID(customerID string) ([]*models.Order, error) { return nil, nil }
func (f *fakeOrderRepo) UpdateStatus(id string, status models.OrderStatus) error { return nil }
func (f *fakeOrderRepo) Delete(id string) error { return nil }

var _ ports.OrderRepository = (*fakeOrderRepo)(nil)

type fakePredictionRepo struct{}

func (f *fakePredictionRepo) Save(prediction *models.Prediction) error { return nil }
func (f *fakePredictionRepo) FindByServiceType(serviceType string) ([]*models.Prediction, error) { return nil, nil }
func (f *fakePredictionRepo) GetHistoricalData(serviceType string) ([]predictor.DataPoint, error) { return nil, nil }

func TestTimerWorkerOnlyMarksReadyAfterEstimatedTime(t *testing.T) {
	orderRepo := &fakeOrderRepo{orders: []*models.Order{{
		ID:            uuid.NewString(),
		Status:        models.StatusProcessing,
		CreatedAt:     time.Now().Add(-5 * time.Minute),
		EstimatedTime: 10 * time.Minute,
	}}}
	updateOrderStatus := orderUseCases.NewUpdateOrderStatusUseCase(orderRepo, &fakePredictionRepo{})
	worker := NewTimerWorker(orderRepo, updateOrderStatus, &fakeNotifier{})

	worker.checkOrders()
}

type fakeNotifier struct{}

func (f *fakeNotifier) NotifyOrderReady(customerID, orderID string) error { return nil }
func (f *fakeNotifier) NotifyStatusChange(customerID, orderID, status string) error { return nil }
