package models

import (
	"strconv"
	"time"

	"github.com/manyminds/api2go/jsonapi"
)

// User model
type User struct {
	ID          uint         `json:"-"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Gender      string       `json:"gender"`
	ExternalID  string       `json:"externalId"`
	Groups      []Group      `json:"-" gorm:"many2many:group_users;"`
	Friendships []Friendship `json:"-"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

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
