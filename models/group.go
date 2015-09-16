package models

import (
	"strconv"
	"time"

	"payup/config"

	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/speps/go-hashids"
)

// Group model that users wil use
type Group struct {
	ID           uint `jsonapi:"-"`
	Name         string
	Description  string
	Users        []User        `gorm:"many2many:group_users;" jsonapi:"-"`
	Transactions []Transaction `jsonapi:"-"`
	HashID       string        `jsonapi:"name=hashId"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `jsonapi:"-"`
}

// GetID returns a stringified version of an ID
func (g Group) GetID() string {
	return strconv.FormatUint(uint64(g.ID), 10)
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
