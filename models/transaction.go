package models

import "time"

// Transaction model
type Transaction struct {
	ID         uint       `json:"id"`
	Amount     int32      `json:"amount"`
	Comment    string     `json:"comment"`
	LenderID   uint       `json:"lenderId"`
	BurrowerID uint       `json:"burrowerId"`
	GroupID    uint       `json:"groupId"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
}
