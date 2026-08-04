package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/pkg/predictor"
	"testing"
	"time"

	"github.com/google/uuid"
)

type createOrderOrderRepoStub struct {
	created *models.Order
}

func (s *createOrderOrderRepoStub) Create(order *models.Order) error {
	s.created = order
	return nil
}

func (s *createOrderOrderRepoStub) FindByID(id string) (*models.Order, error) { return nil, nil }
func (s *createOrderOrderRepoStub) FindAll() ([]*models.Order, error)         { return nil, nil }
func (s *createOrderOrderRepoStub) FindByUserID(userID uuid.UUID) ([]*models.Order, error) {
	return nil, nil
}
func (s *createOrderOrderRepoStub) FindByCustomerID(customerID string) ([]*models.Order, error) {
	return nil, nil
}
func (s *createOrderOrderRepoStub) UpdateStatus(id string, status models.OrderStatus) error {
	return nil
}
func (s *createOrderOrderRepoStub) Delete(id string) error { return nil }

type createOrderPredRepoStub struct {
	historical []predictor.DataPoint
	saved      *models.Prediction
}

func (s *createOrderPredRepoStub) Save(prediction *models.Prediction) error {
	s.saved = prediction
	return nil
}

func (s *createOrderPredRepoStub) FindByServiceType(serviceType string) ([]*models.Prediction, error) {
	return nil, nil
}
func (s *createOrderPredRepoStub) GetHistoricalData(serviceType string) ([]predictor.DataPoint, error) {
	return s.historical, nil
}

type createOrderServiceTypeRepoStub struct {
	serviceType *models.ServiceType
}

func (s *createOrderServiceTypeRepoStub) Create(serviceType *models.ServiceType) error { return nil }
func (s *createOrderServiceTypeRepoStub) FindAll() ([]*models.ServiceType, error)      { return nil, nil }
func (s *createOrderServiceTypeRepoStub) FindByName(name string) (*models.ServiceType, error) {
	return s.serviceType, nil
}
func (s *createOrderServiceTypeRepoStub) Update(serviceType *models.ServiceType) error { return nil }
func (s *createOrderServiceTypeRepoStub) Delete(name string) error                    { return nil }

func TestCreateOrderUseCaseSavesPredictionAndNormalizesServiceType(t *testing.T) {
	orderRepo := &createOrderOrderRepoStub{}
	predRepo := &createOrderPredRepoStub{}
	serviceTypeRepo := &createOrderServiceTypeRepoStub{serviceType: &models.ServiceType{BasePrice: 20, PricePerWeight: 5, PricePerPiece: 2}}
	uc := NewCreateOrderUseCase(orderRepo, predRepo, serviceTypeRepo)

	order, err := uc.Execute("cust-1", "Lavado y Secado", 3, "", 4.5)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if order == nil {
		t.Fatal("expected order, got nil")
	}
	if predRepo.saved == nil {
		t.Fatal("expected prediction to be saved")
	}
	if predRepo.saved.ServiceType != "lavado_y_secado" {
		t.Fatalf("expected normalized service type lavado_y_secado, got %s", predRepo.saved.ServiceType)
	}
	if order.EstimatedCost != 48.5 {
		t.Fatalf("expected estimated cost 48.50, got %.2f", order.EstimatedCost)
	}
	if predRepo.saved.Weight != 4.5 {
		t.Fatalf("expected saved weight 4.5, got %.2f", predRepo.saved.Weight)
	}
	if predRepo.saved.Estimated != order.EstimatedTime {
		t.Fatalf("expected saved estimated time %v, got %v", order.EstimatedTime, predRepo.saved.Estimated)
	}
	if order.ID == "" || predRepo.saved.ID == "" {
		t.Fatal("expected generated ids")
	}
	if _, err := uuid.Parse(order.ID); err != nil {
		t.Fatalf("expected valid order id, got %q: %v", order.ID, err)
	}
	if _, err := uuid.Parse(predRepo.saved.ID); err != nil {
		t.Fatalf("expected valid prediction id, got %q: %v", predRepo.saved.ID, err)
	}
	if time.Since(predRepo.saved.CreatedAt) > time.Minute {
		t.Fatal("expected recent prediction timestamp")
	}
}
