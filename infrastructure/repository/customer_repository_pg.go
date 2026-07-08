package repository

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"database/sql"
)

type customerRepositoryPG struct {
	db *sql.DB
}

var _ ports.CustomerRepository = (*customerRepositoryPG)(nil)

func NewCustomerRepository(db *sql.DB) ports.CustomerRepository {
	return &customerRepositoryPG{db: db}
}

func (r *customerRepositoryPG) Create(customer *models.Customer) error {
	query := `
		INSERT INTO customers (id, name, phone, email) VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query,
		customer.ID,
		customer.Name,
		customer.Phone,
		customer.Email,
	)
	return err
}

func (r *customerRepositoryPG) FindAll() ([]*models.Customer, error) {
	query := `SELECT id, name, email, phone FROM customers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}
	return customers, nil
}

func (r *customerRepositoryPG) FindByID(customerID string) (*models.Customer, error) {
	query := `SELECT id, name, email, phone FROM customers WHERE id = $1`
	row := r.db.QueryRow(query, customerID)

	var customer models.Customer
	if err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
	); err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepositoryPG) Delete(customerID string) error {
	_, err := r.db.Exec(`DELETE FROM customers WHERE id = $1`, customerID)
	return err
}
