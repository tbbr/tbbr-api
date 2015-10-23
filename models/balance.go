package models

import (
	"strconv"
	"time"
)

// Balance model
type Balance struct {
	ID            uint `json:"id"`
	Amount        int
	UserID        uint
	RelatedUserID uint
	GroupID       uint
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

// GetID returns a stringified version of an ID
func (b Balance) GetID() string {
	return strconv.FormatUint(uint64(b.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (b *Balance) SetID(id string) error {
	balanceID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	b.ID = uint(balanceID)
	return nil
}
