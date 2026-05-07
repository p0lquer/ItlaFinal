// Esta struct implementa el PORT definido en domain/ports
// Go verifica esto en tiempo de compilación con la línea del var _
package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
	"time"
)

type orderRepositoryPG struct {
	db *sql.DB
}

// ✅ Go valida en compilación que implementa la interfaz
var _ ports.OrderRepository = (*orderRepositoryPG)(nil)

func NewOrderRepository(db *sql.DB) ports.OrderRepository {
	return &orderRepositoryPG{db: db}
}

func (r *orderRepositoryPG) Create(order *models.Order) error {
	query := `
        INSERT INTO orders (id, customer_id, service_type, pieces_count, notes, status, estimated_time, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := r.db.Exec(query,
		order.ID,
		order.CustomerID,
		order.ServiceType,
		order.PiecesCount,
		order.Notes,
		order.Status,
		order.EstimatedTime.Minutes(),
		order.CreatedAt,
		order.UpdatedAt,
	)
	return err
}

func (r *orderRepositoryPG) FindByID(id string) (*models.Order, error) {
	query := `SELECT id, customer_id, service_type, pieces_count, notes, status, estimated_time, created_at, updated_at FROM orders WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var order models.Order
	var estimatedMinutes float64

	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.ServiceType,
		&order.PiecesCount,
		&order.Notes,
		&order.Status,
		&estimatedMinutes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	order.EstimatedTime = time.Duration(estimatedMinutes) * time.Minute
	return &order, nil
}

func (r *orderRepositoryPG) FindAll() ([]*models.Order, error) {
	query := `SELECT id, customer_id, service_type, pieces_count, notes, status, estimated_time, created_at, updated_at FROM orders ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		var estimatedMinutes float64
		if err := rows.Scan(&order.ID, &order.CustomerID, &order.ServiceType, &order.PiecesCount, &order.Notes, &order.Status, &estimatedMinutes, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		order.EstimatedTime = time.Duration(estimatedMinutes) * time.Minute
		orders = append(orders, &order)
	}
	return orders, nil
}

func (r *orderRepositoryPG) FindByCustomerID(customerID string) ([]*models.Order, error) {
	query := `SELECT id, customer_id, service_type, pieces_count, notes, status, estimated_time, created_at, updated_at FROM orders WHERE customer_id = $1`
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		var estimatedMinutes float64
		rows.Scan(&order.ID, &order.CustomerID, &order.ServiceType, &order.PiecesCount, &order.Notes, &order.Status, &estimatedMinutes, &order.CreatedAt, &order.UpdatedAt)
		order.EstimatedTime = time.Duration(estimatedMinutes) * time.Minute
		orders = append(orders, &order)
	}
	return orders, nil
}

func (r *orderRepositoryPG) UpdateStatus(id string, status models.OrderStatus) error {
	_, err := r.db.Exec(`UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`, status, time.Now(), id)
	return err
}

func (r *orderRepositoryPG) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM orders WHERE id = $1`, id)
	return err
}
