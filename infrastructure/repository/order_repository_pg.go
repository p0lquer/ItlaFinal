// Esta struct implementa el PORT definido en domain/ports
// Go verifica esto en tiempo de compilación con la línea del var _
package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type orderRepositoryPG struct {
	db *sql.DB
}

// ✅ Go valida en compilación que implementa la interfaz
var _ ports.OrderRepository = (*orderRepositoryPG)(nil)

func NewOrderRepository(db *sql.DB) ports.OrderRepository {
	return &orderRepositoryPG{db: db}
}

const orderColumns = "id, customer_id, service_type, pieces_count, weight, price, notes, status, estimated_time, created_at, updated_at, ready_at"

func (r *orderRepositoryPG) Create(order *models.Order) error {
	query := `
        INSERT INTO orders (` + orderColumns + `)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `
	_, err := r.db.Exec(
		query,
		order.ID,
		order.CustomerID,
		order.ServiceType,
		order.PiecesCount,
		order.Weight,
		order.EstimatedCost,
		order.Notes,
		order.Status,
		order.EstimatedTime.Minutes(),
		order.CreatedAt,
		order.UpdatedAt,
		order.ReadyAt,
	)
	return err
}

func scanOrder(scan func(dest ...any) error) (*models.Order, error) {
	var order models.Order
	var estimatedMinutes float64
	var weight, price sql.NullFloat64
	var readyAt sql.NullTime

	err := scan(
		&order.ID,
		&order.CustomerID,
		&order.ServiceType,
		&order.PiecesCount,
		&weight,
		&price,
		&order.Notes,
		&order.Status,
		&estimatedMinutes,
		&order.CreatedAt,
		&order.UpdatedAt,
		&readyAt,
	)

	if err != nil {
		return nil, err
	}

	if weight.Valid {
		order.Weight = weight.Float64
	}

	if price.Valid {
		order.EstimatedCost = price.Float64
	}

	if readyAt.Valid {
		order.ReadyAt = &readyAt.Time
	}

	order.EstimatedTime = time.Duration(estimatedMinutes) * time.Minute

	return &order, nil
}

func (r *orderRepositoryPG) FindByID(id string) (*models.Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var order models.Order
	var estimatedMinutes float64
	var readyAt sql.NullTime
	var weight sql.NullFloat64
	var price sql.NullFloat64

	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.ServiceType,
		&order.PiecesCount,
		&weight,
		&price,
		&order.Notes,
		&order.Status,
		&estimatedMinutes,
		&order.CreatedAt,
		&order.UpdatedAt,
		&readyAt,
	)
	if err != nil {
		return nil, err
	}

	if weight.Valid {
		order.Weight = weight.Float64
	}

	if price.Valid {
		order.EstimatedCost = price.Float64
	}

	if readyAt.Valid {
		order.ReadyAt = &readyAt.Time
	}

	order.EstimatedTime = time.Duration(estimatedMinutes) * time.Minute

	return &order, nil
}

func (r *orderRepositoryPG) FindAll() ([]*models.Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order, err := scanOrder(rows.Scan)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *orderRepositoryPG) FindByCustomerID(customerID string) ([]*models.Order, error) {
	rows, err := r.db.Query(`SELECT `+orderColumns+` FROM orders WHERE customer_id = $1 ORDER BY created_at DESC`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order, err := scanOrder(rows.Scan)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *orderRepositoryPG) UpdateStatus(id string, status models.OrderStatus) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE orders 
		 SET status = $1, 
		     updated_at = $2, 
		     ready_at = CASE WHEN $1::varchar = 'lista' THEN $2 ELSE ready_at END 
		 WHERE id = $3`,
		status, now, id,
	)
	return err
}
func (r *orderRepositoryPG) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM orders WHERE id = $1`, id)
	return err
}
func (r *orderRepositoryPG) FindByUserID(userID uuid.UUID) ([]*models.Order, error) {
	query := `SELECT ` + orderColumns + `
          FROM orders
          WHERE customer_id = $1
          ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []*models.Order

	for rows.Next() {
		order, err := scanOrder(rows.Scan)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, rows.Err()
}
