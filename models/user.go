package models

import (
	"time"
)

// User model
type User struct {
	ID        int
	Name      string
	Username  string
	Email     string
	Groups    []Group `gorm:"many2many:group_users;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
