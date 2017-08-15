package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/tbbr/tbbr-api/hashid"
)

// Friendship model
type Friendship struct {
	ID               uint       `json:"-"`
	UserID           uint       `json:"userId"`
	FriendID         uint       `json:"friendId"`
	FriendshipDataID uint       `json:"friendshipDataId"`
	HashID           string     `json:"hashId"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"-"`
	User             User       `json:"-" sql:"-"`
	Friend           User       `json:"-" sql:"-"`
	Balance          int        `json:"balance" sql:"-"`
}

// BeforeCreate will Find or Create FriendshipData model
func (f *Friendship) BeforeCreate(db *gorm.DB) (err error) {
	// Try to find FriendshipData
	var otherFriendship Friendship
	if db.Where("user_id = ? AND friend_id = ?", f.FriendID, f.UserID).First(&otherFriendship).RecordNotFound() {
		// Create FriendshipData
		var fd FriendshipData
		fd.Balance = 0
		fd.PositiveUserID = f.UserID
		db.Create(&fd)
		f.FriendshipDataID = fd.ID
	} else {
		f.FriendshipDataID = otherFriendship.FriendshipDataID
	}
	return
}

// AfterCreate generates a HashID for a Friendship based on it's numeric ID field
func (f *Friendship) AfterCreate(db *gorm.DB) (err error) {
	f.HashID = hashid.Generate(f.ID)
	db.Save(&f)
	return
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

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

// GetReferencedStructs to satisfy the jsonapi.MarshalIncludedRelations interface
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
