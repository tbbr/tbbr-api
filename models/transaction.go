package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/manyminds/api2go/jsonapi"
)

// Transaction model
type Transaction struct {
	ID            uint `json:"id"`
	Type          string
	Amount        int
	Memo          string
	RelatedUserID uint `jsonapi:"name=relatedUserId"`
	GroupID       uint `jsonapi:"name=groupId"`
	BalanceID     uint `jsonapi:"name=balanceId"`
	CreatorID     uint `jsonapi:"name=creatorId"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time

	Creator     User `jsonapi:"-"`
	RelatedUser User `jsonapi:"-"`
}

// GetID returns a stringified version of an ID
func (t Transaction) GetID() string {
	return strconv.FormatUint(uint64(t.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (t *Transaction) SetID(id string) error {
	transactionID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	t.ID = uint(transactionID)
	return nil
}

// GetReferences returns all related structs to groups
func (t Transaction) GetReferences() []jsonapi.Reference {
	// TODO:: Uncommenting the groups relation will endup with an empty array
	// relation for users in the response, which isn't necessarily true
	// We'll need to fix this on the routeHandler level

	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "creator",
		},
		{
			Type: "users",
			Name: "related-user",
		},
	}
}

// GetReferencedIDs satisfies the jsonapi.MarshalLinkedRelations interface
func (t Transaction) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	result = append(result, jsonapi.ReferenceID{
		ID:   strconv.FormatUint(uint64(t.CreatorID), 10),
		Type: "users",
		Name: "creator",
	})

	result = append(result, jsonapi.ReferenceID{
		ID:   strconv.FormatUint(uint64(t.RelatedUserID), 10),
		Type: "users",
		Name: "related-user",
	})
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (t Transaction) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	result = append(result, t.Creator)
	result = append(result, t.RelatedUser)

	return result
}

// SetToOneReferenceID sets the reference ID and satisfies the jsonapi.UnmarshalToOneRelations interface
func (t *Transaction) SetToOneReferenceID(name, ID string) error {
	temp, err := strconv.ParseUint(ID, 10, 64)

	if err != nil {
		return err
	}

	switch name {
	case "related-user":
		t.RelatedUserID = uint(temp)
	case "group-id":
		t.GroupID = uint(temp)
	case "balance-id":
		t.BalanceID = uint(temp)
	case "creator":
		t.CreatorID = uint(temp)
	}

	return errors.New("There is no to-one relationship with the name " + name)
}
