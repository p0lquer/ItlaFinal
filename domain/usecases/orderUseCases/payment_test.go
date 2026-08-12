package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"testing"
)

type paymentRepoStub struct {
	payment *models.Payment
	created *models.Payment
}

func (r *paymentRepoStub) FindByOrderID(string) (*models.Payment, error) { return r.payment, nil }
func (r *paymentRepoStub) FindByCustomerID(string) ([]*models.Payment, error) {
	return nil, nil
}
func (r *paymentRepoStub) Create(payment *models.Payment) error { r.created = payment; return nil }

var _ ports.PaymentRepository = (*paymentRepoStub)(nil)

func TestPaymentIsOnlyAvailableForReadyOrder(t *testing.T) {
	for _, status := range []models.OrderStatus{models.StatusReceived, models.StatusProcessing, models.StatusDelivered} {
		uc := NewPaymentUseCase(&statusOrderRepo{order: &models.Order{ID: "order-1", Status: status, EstimatedCost: 120}}, &paymentRepoStub{})
		summary, err := uc.Summary("order-1")
		if err != nil || summary.CanPay {
			t.Fatalf("status %q must not be payable: summary=%#v err=%v", status, summary, err)
		}
		if _, err := uc.Pay("order-1", "card"); err == nil {
			t.Fatalf("status %q accepted a payment", status)
		}
	}
	uc := NewPaymentUseCase(&statusOrderRepo{order: &models.Order{ID: "order-2", Status: models.StatusReady, EstimatedCost: 120}}, &paymentRepoStub{})
	summary, err := uc.Summary("order-2")
	if err != nil || !summary.CanPay {
		t.Fatalf("ready order must be payable: summary=%#v err=%v", summary, err)
	}
}
