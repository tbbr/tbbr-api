package models

import "strconv"

// Balance model
type Balance struct {
	ID             uint
	Amount         int
	PositiveUserID uint
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
