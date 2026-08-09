// Esta struct implementa el PORT definido en domain/ports
// Go verifica esto en tiempo de compilación con la línea del var _
package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
	"fmt"
	"strings"
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
// transitions describe de qué estados de origen se permite pasar a target.
// Solo se avanza de forma adyacente/semádyacente para que no existan saltos
// que dejen la orden sin registro ni dobles marcas de "lista".
var transitions = map[models.OrderStatus][]models.OrderStatus{
	models.StatusReceived:   {},
	models.StatusProcessing: {models.StatusReceived},
	models.StatusReady:      {models.StatusProcessing, models.StatusReceived},
	models.StatusDelivered:  {models.StatusProcessing, models.StatusReady},
}
// Transition es la transición guardada e idempotente.
// Devuelve (true, nil) solo cuando cambió el estado de la fila; false si la
// orden ya estaba en target o no hay transición válida desde su estado actual.
// También fija started_at al entrar a "en_proceso" y ready_at al llegar a "lista".
func (r *orderRepositoryPG) Transition(orderID string, target models.OrderStatus) (bool, error) {
	from, ok := transitions[target]
	if !ok || len(from) == 0 {
		return false, nil
	}
	placeholders := make([]string, len(from))
	args := make([]any, 0, len(from)+2)
	args = append(args, orderID, string(target))
	for i, s := range from {
		placeholders[i] = fmt.Sprintf("$%d", i+3)
		args = append(args, string(s))
	}

	res, err := r.db.Exec(fmt.Sprintf(`
		UPDATE orders
		   SET status = $2,
		       updated_at = now(),
		       started_at = CASE WHEN $2 = 'en_proceso' AND started_at IS NULL THEN now() ELSE started_at END,
		       ready_at   = CASE WHEN $2 = 'lista'      AND ready_at   IS NULL THEN now() ELSE ready_at   END
		 WHERE id = $1
		   AND status <> $2
		   AND status IN (%s)`, strings.Join(placeholders, ", ")), args...)
	if err != nil {
		return false, fmt.Errorf("transition order: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("transition order rows: %w", err)
	}
	return n > 0, nil
}