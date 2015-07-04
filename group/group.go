package group

import (
	"payup/user"
)

// Group model that users wil use
type Group struct {
	ID          int
	Name        string
	Description string
	Users       []user.User `gorm:"many2many:group_users;"`
}
