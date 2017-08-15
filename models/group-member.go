package models

import (
	"strconv"
	"time"
)

// GroupMember model keeps track of user's membership of a group
// and the amount they've send, and received
type GroupMember struct {
	ID             uint       `json:"-"`
	GroupID        uint       `json:"groupId"`
	UserID         uint       `json:"userId"`
	AmountSent     uint       `json:"amountSent"`
	AmountReceived uint       `json:"amountReceived"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"-"`

	Group Group `json:"-" sql:"-"`
	User  User  `json:"-" sql:"-"`
}

func (gm GroupMember) SetDefault() {
	gm.AmountSent = 0
	gm.AmountReceived = 0
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

// GetID returns a stringified version of an ID
func (gm GroupMember) GetID() string {
	return strconv.FormatUint(uint64(gm.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (gm *GroupMember) SetID(id string) error {
	groupMemberID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	gm.ID = uint(groupMemberID)
	return nil
}
