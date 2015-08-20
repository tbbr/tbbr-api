package models

import "time"

// Group model that users wil use
type Group struct {
	ID          int
	Name        string
	Description string
	Users       []User `gorm:"many2many:group_users;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
