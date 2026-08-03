package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"testing"
	"time"
)

type stubOrderRepository struct {
	orders []*models.Order
}

func (s *stubOrderRepository) Create(order *models.Order) error { return nil }
func (s *stubOrderRepository) FindByID(id string) (*models.Order, error) {
	return nil, nil
}
func (s *stubOrderRepository) FindAll() ([]*models.Order, error) { return nil, nil }
func (s *stubOrderRepository) FindByUserID(userID string) ([]*models.Order, error) {
	return nil, nil
}
func (s *stubOrderRepository) FindByCustomerID(customerID string) ([]*models.Order, error) {
	return s.orders, nil
}
func (s *stubOrderRepository) UpdateStatus(id string, status models.OrderStatus) error {
	return nil
}
func (s *stubOrderRepository) Delete(id string) error { return nil }

var _ ports.OrderRepository = (*stubOrderRepository)(nil)

func TestGetMyOrdersUseCaseExecuteReturnsOrdersForCustomer(t *testing.T) {
	expectedOrder := &models.Order{
		ID:          "ord-1",
		CustomerID:  "cust-1",
		ServiceType: "lavado_secado",
		CreatedAt:   time.Now(),
	}

	repo := &stubOrderRepository{orders: []*models.Order{expectedOrder}}
	uc := NewGetMyOrdersUseCase(repo)

	got, err := uc.Execute("cust-1")
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 order, got %d", len(got))
	}
	if got[0].ID != expectedOrder.ID {
		t.Fatalf("expected order ID %s, got %s", expectedOrder.ID, got[0].ID)
	}
}
