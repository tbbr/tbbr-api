package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/tbbr/tbbr-api/config"

	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/speps/go-hashids"
)

// Group model that users wil use
type Group struct {
	ID           uint          `json:"-"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Users        []User        `json:"-" gorm:"many2many:group_users;"`
	Transactions []Transaction `json:"-"`
	HashID       string        `json:"hashId"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	UserIDs      []uint        `json:"-" sql:"-"`
}

// AfterCreate generates a HashID for a Group based on it's numeric ID field
func (g *Group) AfterCreate(db *gorm.DB) (err error) {
	hd := hashids.NewData()
	hd.Salt = config.HashID.Salt
	hd.MinLength = config.HashID.MinLength
	h := hashids.NewWithData(hd)

	a := []int{0}
	a[0] = int(g.ID)

	// Encode
	e, _ := h.Encode(a)
	g.HashID = e

	// Save
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

// SetToManyReferenceIDs sets the users reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (g *Group) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "users" {
		for _, i := range IDs {
			j, err := strconv.ParseUint(i, 10, 64)
			if err != nil {
				return err
			}
			g.UserIDs = append(g.UserIDs, uint(j))
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// GetReferences returns all related structs to groups
func (g Group) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "users",
		},
		{
			Type: "transactions",
			Name: "transactions",
		},
	}
}

// GetReferencedIDs satisfies the jsonapi.MarshalLinkedRelations interface
func (g Group) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, user := range g.Users {
		result = append(result, jsonapi.ReferenceID{
			ID:   user.GetID(),
			Type: "users",
			Name: "users",
		})
	}

	for _, transaction := range g.Transactions {
		result = append(result, jsonapi.ReferenceID{
			ID:   transaction.GetID(),
			Type: "transactions",
			Name: "transactions",
		})
	}
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (g Group) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range g.Users {
		result = append(result, g.Users[key])
	}

	return result
}
