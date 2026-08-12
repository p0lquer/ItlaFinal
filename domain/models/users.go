package models

import "time"

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleOperator UserRole = "operator"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string // guardado como hash bcrypt
	Role      UserRole
	IsActive  bool
	CreatedAt time.Time
}

type UserFilter struct {
	Search   string
	Role     UserRole
	IsActive *bool
	Page     int
	PageSize int
}

type UserPage struct {
	Users    []*User
	Total    int
	Page     int
	PageSize int
}
