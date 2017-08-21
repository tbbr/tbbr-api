package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/manyminds/api2go/jsonapi"
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

	Group Group `json:"-"`
	User  User  `json:"-"`
}

func (gm GroupMember) SetDefault() {
	gm.AmountSent = 0
	gm.AmountReceived = 0
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

func (gm GroupMember) GetName() string {
	return "group-members"
}

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

// GetReferences returns all related structs to groupTransactions
func (gm GroupMember) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "user",
		},
		{
			Type: "groups",
			Name: "group",
		},
	}
}

// GetReferencedIDs satisfies the jsonapi.MarshalLinkedRelations interface
func (gm GroupMember) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	result = append(result, jsonapi.ReferenceID{
		ID:   strconv.FormatUint(uint64(gm.UserID), 10),
		Type: "users",
		Name: "user",
	})

	result = append(result, jsonapi.ReferenceID{
		ID:   strconv.FormatUint(uint64(gm.GroupID), 10),
		Type: "groups",
		Name: "group",
	})
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (gm GroupMember) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	result = append(result, gm.User)
	return result
}

// SetToOneReferenceID sets the reference ID and satisfies the jsonapi.UnmarshalToOneRelations interface
func (gm *GroupMember) SetToOneReferenceID(name, ID string) error {
	temp, err := strconv.ParseUint(ID, 10, 64)

	if err != nil {
		return err
	}

	switch name {
	case "user":
		gm.UserID = uint(temp)
		return nil
	case "group":
		gm.GroupID = uint(temp)
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}
