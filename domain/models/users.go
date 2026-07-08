package models

import "time"

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleOperator UserRole = "operator"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string // guardado como hash bcrypt
	Role      UserRole
	CreatedAt time.Time
}
