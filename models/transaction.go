package models

import (
	"strconv"
	"time"
)

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

// GetID returns a stringified version of an ID
func (t Transaction) GetID() string {
	return strconv.FormatUint(uint64(t.ID), 10)
}
