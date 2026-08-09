package ports

// Notifier permite desacoplar el aviso al cliente del worker.
type Notifier interface {
	NotifyOrderReady(customerID string, orderID string) error
	NotifyStatusChange(customerID string, orderID string, status string) error
}
