package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/manyminds/api2go/jsonapi"
)

// Friendship model
type Friendship struct {
	ID        uint `json:"id"`
	UserID    uint `jsonapi:"name=userId"`
	FriendID  uint `jsonapi:"name=friendId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `jsonapi:"-"`

	User   User `jsonapi:"-" sql:"-"`
	Friend User `jsonapi:"-" sql:"-"`
}

// GetID returns a stringified version of an ID
func (f Friendship) GetID() string {
	return strconv.FormatUint(uint64(f.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (f *Friendship) SetID(id string) error {
	friendshipID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	f.ID = uint(friendshipID)
	return nil
}

// GetReferences returns all related structs to friendships
func (f Friendship) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "user",
		},
		{
			Type: "users",
			Name: "friend",
		},
	}
}

// GetReferencedIDs satisfies the jsonapi.MarshalLinkedRelations interface
func (f Friendship) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	result = append(result, jsonapi.ReferenceID{
		ID:   strconv.FormatUint(uint64(f.UserID), 10),
		Type: "users",
		Name: "user",
	})

	result = append(result, jsonapi.ReferenceID{
		ID:   strconv.FormatUint(uint64(f.FriendID), 10),
		Type: "users",
		Name: "friend",
	})
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (f Friendship) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	result = append(result, f.User)
	result = append(result, f.Friend)

	return result
}

// SetToOneReferenceID sets the reference ID and satisfies the jsonapi.UnmarshalToOneRelations interface
func (f *Friendship) SetToOneReferenceID(name, ID string) error {
	temp, err := strconv.ParseUint(ID, 10, 64)

	if err != nil {
		return err
	}

	switch name {
	case "user":
		f.UserID = uint(temp)
	case "friend":
		f.FriendID = uint(temp)
	}

	return errors.New("There is no to-one relationship with the name " + name)
}
