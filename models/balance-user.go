package models

import "strconv"

// BalanceUser model
type BalanceUser struct {
	ID            uint
	BalanceID     uint `jsonapi:"name=balanceId"`
	UserID        uint `jsonapi:"name=userId"`
	RelatedUserID uint `jsonapi:"name=relatedUserId"`
	GroupID       uint `jsonapi:"name=groupId"`
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

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
