package models

import (
	"time"
)

// User model
type User struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Groups    []Group    `gorm:"many2many:group_users;" json:"groups"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
