package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

type orderRepositoryPG struct{ db *sql.DB }

var _ ports.OrderRepository = (*orderRepositoryPG)(nil)

func NewOrderRepository(db *sql.DB) ports.OrderRepository { return &orderRepositoryPG{db: db} }

// Keep this list and scanOrder together. estimated_time is always persisted in minutes.
const orderColumns = "id, customer_id, service_type, pieces_count, weight, estimated_cost, notes, status, estimated_time, created_at, updated_at, ready_at, started_at"

func (r *orderRepositoryPG) Create(order *models.Order) error {
	_, err := r.db.Exec(`INSERT INTO orders (`+orderColumns+`)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		order.ID, order.CustomerID, order.ServiceType, order.PiecesCount, order.Weight,
		order.EstimatedCost, order.Notes, order.Status, order.EstimatedTime.Minutes(),
		order.CreatedAt, order.UpdatedAt, order.ReadyAt, order.StartedAt)
	return err
}

func scanOrder(scan func(...any) error) (*models.Order, error) {
	var order models.Order
	var minutes float64
	var weight, cost sql.NullFloat64
	var readyAt, startedAt sql.NullTime
	if err := scan(&order.ID, &order.CustomerID, &order.ServiceType, &order.PiecesCount,
		&weight, &cost, &order.Notes, &order.Status, &minutes, &order.CreatedAt,
		&order.UpdatedAt, &readyAt, &startedAt); err != nil {
		return nil, err
	}
	if weight.Valid {
		order.Weight = weight.Float64
	}
	if cost.Valid {
		order.EstimatedCost = cost.Float64
	}
	if readyAt.Valid {
		order.ReadyAt = &readyAt.Time
	}
	if startedAt.Valid {
		order.StartedAt = &startedAt.Time
	}
	order.EstimatedTime = time.Duration(minutes * float64(time.Minute))
	return &order, nil
}

func (r *orderRepositoryPG) FindByID(id string) (*models.Order, error) {
	return scanOrder(r.db.QueryRow(`SELECT `+orderColumns+` FROM orders WHERE id=$1`, id).Scan)
}

func findOrders(rows *sql.Rows) ([]*models.Order, error) {
	defer rows.Close()
	orders := make([]*models.Order, 0)
	for rows.Next() {
		order, err := scanOrder(rows.Scan)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *orderRepositoryPG) FindAll() ([]*models.Order, error) {
	rows, err := r.db.Query(`SELECT ` + orderColumns + ` FROM orders ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	return findOrders(rows)
}

// List retrieves a stable, paginated operator view. Search covers the order
// ID, customer name/email and service type, while every dynamic value remains
// a bound SQL parameter.
func (r *orderRepositoryPG) List(filter models.OrderFilter) (*models.OrderPage, error) {
	filter = normalizeOrderFilter(filter)
	where, args := buildOrderListWhere(filter)
	from := ` FROM orders o LEFT JOIN customers c ON c.id = o.customer_id `

	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*)`+from+where, args...).Scan(&total); err != nil {
		return nil, err
	}

	limitPosition := len(args) + 1
	offsetPosition := len(args) + 2
	query := `SELECT ` + prefixedOrderColumns("o") + from + where +
		` ORDER BY o.created_at DESC, o.id DESC LIMIT $` + itoa(limitPosition) + ` OFFSET $` + itoa(offsetPosition)
	rows, err := r.db.Query(query, append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)...)
	if err != nil {
		return nil, err
	}
	orders, err := findOrders(rows)
	if err != nil {
		return nil, err
	}
	return &models.OrderPage{Orders: orders, Total: total, Page: filter.Page, PageSize: filter.PageSize}, nil
}

func normalizeOrderFilter(filter models.OrderFilter) models.OrderFilter {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	filter.Search = strings.TrimSpace(filter.Search)
	filter.CustomerID = strings.TrimSpace(filter.CustomerID)
	filter.ServiceType = strings.TrimSpace(filter.ServiceType)
	return filter
}

func buildOrderListWhere(filter models.OrderFilter) (string, []any) {
	clauses := make([]string, 0, 6)
	args := make([]any, 0, 6)
	add := func(clause string, value any) {
		args = append(args, value)
		clauses = append(clauses, clause+` $`+itoa(len(args)))
	}
	if filter.Search != "" {
		args = append(args, "%"+strings.ToLower(filter.Search)+"%")
		position := itoa(len(args))
		clauses = append(clauses, `(LOWER(o.id) LIKE $`+position+` OR LOWER(o.service_type) LIKE $`+position+` OR LOWER(COALESCE(c.name, '')) LIKE $`+position+` OR LOWER(COALESCE(c.email, '')) LIKE $`+position+`)`)
	}
	if filter.CustomerID != "" {
		add(`o.customer_id =`, filter.CustomerID)
	}
	if filter.Status != "" {
		add(`o.status =`, string(filter.Status))
	}
	if filter.ServiceType != "" {
		add(`LOWER(o.service_type) = LOWER(`, filter.ServiceType)
		clauses[len(clauses)-1] += `)`
	}
	if filter.From != nil {
		add(`o.created_at >=`, *filter.From)
	}
	if filter.To != nil {
		add(`o.created_at <=`, *filter.To)
	}
	if len(clauses) == 0 {
		return "", args
	}
	return ` WHERE ` + strings.Join(clauses, ` AND `), args
}

func prefixedOrderColumns(prefix string) string {
	return prefix + `.` + strings.ReplaceAll(orderColumns, `, `, `, `+prefix+`.`)
}

// strconv.Itoa is avoided here because it is only used to build placeholder
// positions that are calculated locally, never from client input.
func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	result := make([]byte, 0, 3)
	for value > 0 {
		result = append([]byte{byte('0' + value%10)}, result...)
		value /= 10
	}
	return string(result)
}

func (r *orderRepositoryPG) FindByCustomerID(customerID string) ([]*models.Order, error) {
	rows, err := r.db.Query(`SELECT `+orderColumns+` FROM orders WHERE customer_id=$1 ORDER BY created_at DESC`, customerID)
	if err != nil {
		return nil, err
	}
	return findOrders(rows)
}

func (r *orderRepositoryPG) FindStatusHistory(orderID string) ([]*models.OrderStatusChange, error) {
	rows, err := r.db.Query(`SELECT order_id, from_status, to_status, changed_by, changed_by_role, description, changed_at
		FROM order_status_history WHERE order_id=$1 ORDER BY changed_at ASC`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]*models.OrderStatusChange, 0)
	for rows.Next() {
		change := &models.OrderStatusChange{}
		var changedBy, changedByRole, description sql.NullString
		if err := rows.Scan(&change.OrderID, &change.FromStatus, &change.ToStatus, &changedBy, &changedByRole, &description, &change.ChangedAt); err != nil {
			return nil, err
		}
		if changedBy.Valid {
			change.ChangedBy = changedBy.String
		}
		if changedByRole.Valid {
			change.ChangedByRole = changedByRole.String
		}
		if description.Valid {
			change.Description = description.String
		}
		history = append(history, change)
	}
	return history, rows.Err()
}

// FindByUserID uses the explicit customer/user relationship; customer IDs are
// generated by the server and are intentionally no longer assumed to equal user IDs.
func (r *orderRepositoryPG) FindByUserID(userID uuid.UUID) ([]*models.Order, error) {
	rows, err := r.db.Query(`SELECT o.`+strings.ReplaceAll(orderColumns, ", ", ", o.")+`
		FROM orders o JOIN customers c ON c.id=o.customer_id
		WHERE c.user_id=$1 ORDER BY o.created_at DESC`, userID.String())
	if err != nil {
		return nil, err
	}
	return findOrders(rows)
}

func (r *orderRepositoryPG) UpdateStatus(id string, status models.OrderStatus) error {
	_, err := r.db.Exec(`UPDATE orders SET status=$1, updated_at=now(),
		ready_at=CASE WHEN $1='lista' AND ready_at IS NULL THEN now() ELSE ready_at END WHERE id=$2`, status, id)
	return err
}

func (r *orderRepositoryPG) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM orders WHERE id=$1`, id)
	return err
}

var transitions = map[models.OrderStatus]models.OrderStatus{
	models.StatusProcessing: models.StatusReceived,
	models.StatusReady:      models.StatusProcessing,
	models.StatusDelivered:  models.StatusReady,
}

func (r *orderRepositoryPG) Transition(orderID string, target models.OrderStatus, change models.OrderStatusChange) (bool, error) {
	from, ok := transitions[target]
	if !ok {
		return false, nil
	}
	if change.ChangedBy == "" {
		change.ChangedBy = "sistema"
	}
	if change.ChangedByRole == "" {
		change.ChangedByRole = "system"
	}

	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`UPDATE orders SET status=$2, updated_at=now(),
		started_at=CASE WHEN $2='en_proceso' AND started_at IS NULL THEN now() ELSE started_at END,
		ready_at=CASE WHEN $2='lista' AND ready_at IS NULL THEN now() ELSE ready_at END
		WHERE id=$1 AND status=$3`, orderID, string(target), string(from))
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	if err != nil || n == 0 {
		return n > 0, err
	}

	_, err = tx.Exec(`INSERT INTO order_status_history
		(order_id, from_status, to_status, changed_by, changed_by_role, description)
		VALUES ($1, $2, $3, $4, $5, $6)`, orderID, string(from), string(target), change.ChangedBy, change.ChangedByRole, change.Description)
	if err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}
