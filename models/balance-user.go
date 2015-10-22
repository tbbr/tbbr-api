package models

import (
	"strconv"
	"time"
)

// BalanceUser model
type BalanceUser struct {
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
func (bu BalanceUser) GetID() string {
	return strconv.FormatUint(uint64(bu.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (bu *BalanceUser) SetID(id string) error {
	balanceUserID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	bu.ID = uint(balanceUserID)
	return nil
}
