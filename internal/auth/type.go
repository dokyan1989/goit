package auth

import (
	"database/sql"
	"time"
)

type User struct {
	UserID      int64
	UserName    string
	Password    string
	Email       string
	FirstName   string
	LastName    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLoginAt sql.NullTime
	Status      string
}

type UserWithRoles struct {
	User
	Roles []Role
}

type Role struct {
	RoleID      int64
	RoleName    string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string
}

type UserRole struct {
	UserRoleID int64
	UserID     int64
	RoleID     int64
}
