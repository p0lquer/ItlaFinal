package ports

// La implementación real estará en infrastructure/websocket
type Notifier interface {
	NotifyOrderReady(customerID string, orderID string) error
	NotifyStatusChange(customerID string, orderID string, status string) error
}
