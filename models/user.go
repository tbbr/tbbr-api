package models

import (
	"strconv"
	"time"

	"github.com/manyminds/api2go/jsonapi"
)

// User model
type User struct {
	ID           uint `jsonapi:"-"`
	Name         string
	Email        string
	Gender       string
	AvatarURL    string        `jsonapi:"name=avatarUrl"`
	ExternalID   string        `jsonapi:"-"`
	Groups       []Group       `gorm:"many2many:group_users;" jsonapi:"-"`
	Friendships  []Friendship  `jsonapi:"-"`
	BalanceUsers []BalanceUser `jsonapi:"-"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GetID returns a stringified version of an ID
func (u User) GetID() string {
	return strconv.FormatUint(uint64(u.ID), 10)
}

// GetReferences returns all related structs to groups
func (u User) GetReferences() []jsonapi.Reference {
	// TODO:: Uncommenting the groups relation will endup with an empty array
	// relation for users in the response, which isn't necessarily true
	// We'll need to fix this on the routeHandler level

	return []jsonapi.Reference{
	// {
	// 	Type: "groups",
	// 	Name: "groups",
	// },
	}
}

// GetReferencedIDs satisfies the jsonapi.MarshalLinkedRelations interface
func (u User) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	// for _, group := range u.Groups {
	// 	result = append(result, jsonapi.ReferenceID{
	// 		ID:   group.GetID(),
	// 		Type: "groups",
	// 		Name: "groups",
	// 	})
	// }
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u User) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	// for key := range u.Groups {
	// 	result = append(result, u.Groups[key])
	// }

	return result
}
