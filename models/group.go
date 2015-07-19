package models

// Group model that users wil use
type Group struct {
	ID          int
	Name        string
	Description string
	Users       []User `gorm:"many2many:group_users;"`
}
