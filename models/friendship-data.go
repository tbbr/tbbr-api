package models

import "strconv"

// FriendshipData model
type FriendshipData struct {
	ID             uint
	Balance        int
	PositiveUserID uint
}

// TableName gives gorm information on the name of the table
func (fd FriendshipData) TableName() string {
	return "friendship_data"
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

// GetID returns a stringified version of an ID
func (fd FriendshipData) GetID() string {
	return strconv.FormatUint(uint64(fd.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (fd *FriendshipData) SetID(id string) error {
	friendshipDataID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	fd.ID = uint(friendshipDataID)
	return nil
}
