package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/tbbr/tbbr-api/hashid"
)

// Group model that users wil use
type Group struct {
	ID          uint      `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	HashID      string    `json:"hashId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	GroupMemberIDs []uint        `json:"-" sql:"-"`
	GroupMembers   []GroupMember `json:"-"`
}

// AfterCreate generates a HashID for a Group based on it's numeric ID field
func (g *Group) AfterCreate(db *gorm.DB) (err error) {
	g.HashID = hashid.Generate(g.ID)
	db.Save(&g)
	return
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

// GetID returns a stringified version of an ID
func (g Group) GetID() string {
	return strconv.FormatUint(uint64(g.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (g *Group) SetID(id string) error {
	groupID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	g.ID = uint(groupID)
	return nil
}

// SetToManyReferenceIDs sets the groupMembers reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (g *Group) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "group-members" {
		for _, i := range IDs {
			j, err := strconv.ParseUint(i, 10, 64)
			if err != nil {
				return err
			}
			g.GroupMemberIDs = append(g.GroupMemberIDs, uint(j))
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// GetReferences returns all related structs to groups
func (g Group) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "group-members",
			Name: "groupMembers",
		},
	}
}

// GetReferencedIDs satisfies the jsonapi.MarshalLinkedRelations interface
func (g Group) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, member := range g.GroupMembers {
		result = append(result, jsonapi.ReferenceID{
			ID:   member.GetID(),
			Type: "group-members",
			Name: "groupMembers",
		})
	}
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (g Group) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range g.GroupMembers {
		result = append(result, g.GroupMembers[key])
	}

	return result
}
