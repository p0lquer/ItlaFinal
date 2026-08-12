package orderUseCases

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PaymentUseCase struct {
	orders   ports.OrderRepository
	payments ports.PaymentRepository
}

func NewPaymentUseCase(orders ports.OrderRepository, payments ports.PaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{orders, payments}
}

func (uc *PaymentUseCase) Summary(orderID string) (*models.PaymentSummary, error) {
	o, err := uc.orders.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	p, err := uc.payments.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	// Payment is only accepted once production has explicitly marked the order
	// ready for pickup. UI rules are helpful, but this domain rule is final.
	return &models.PaymentSummary{Order: o, Payment: p, CanPay: p == nil && o.Status == models.StatusReady}, nil
}

func (uc *PaymentUseCase) History(customerID string) ([]*models.Payment, error) {
	if strings.TrimSpace(customerID) == "" {
		return []*models.Payment{}, nil
	}
	return uc.payments.FindByCustomerID(customerID)
}

func (uc *PaymentUseCase) Pay(orderID, method string) (*models.Payment, error) {
	method = strings.ToLower(strings.TrimSpace(method))
	switch method {
	case "cash":
		method = "efectivo"
	case "card":
		method = "tarjeta"
	case "transfer":
		method = "transferencia"
	case "efectivo", "tarjeta", "transferencia":
	default:
		method = ""
	}
	if method == "" {
		return nil, errors.New("metodo de pago invalido")
	}
	s, err := uc.Summary(orderID)
	if err != nil {
		return nil, err
	}
	if s.Payment != nil {
		return nil, errors.New("esta orden ya fue pagada")
	}
	if !s.CanPay {
		return nil, errors.New("la orden aun no esta disponible para pago")
	}
	now := time.Now()
	p := &models.Payment{ID: uuid.NewString(), OrderID: orderID, Amount: s.Order.EstimatedCost, Currency: "DOP", Method: method, Status: "paid", ReceiptNo: "TGB-" + now.Format("20060102") + "-" + strings.ToUpper(uuid.NewString()[:8]), PaidAt: now}
	if err := uc.payments.Create(p); err != nil {
		// The database unique constraint is authoritative when simultaneous requests race.
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, errors.New("esta orden ya fue pagada")
		}
		return nil, err
	}
	return p, nil
}
