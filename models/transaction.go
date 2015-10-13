package models

import (
	"errors"
	"strconv"
	"time"
)

// Transaction model
type Transaction struct {
	ID         uint `json:"id"`
	Type       string
	Amount     int        `json:"amount"`
	Comment    string     `json:"comment"`
	LenderID   uint       `json:"lenderId"`
	BurrowerID uint       `json:"burrowerId"`
	GroupID    uint       `json:"groupId"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
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

// SetToOneReferenceID sets the reference ID and satisfies the jsonapi.UnmarshalToOneRelations interface
func (t *Transaction) SetToOneReferenceID(name, ID string) error {
	temp, err := strconv.ParseUint(ID, 10, 64)

	if err != nil {
		return err
	}

	switch name {
	case "lender-id":
		t.LenderID = uint(temp)
	case "burrower-id":
		t.BurrowerID = uint(temp)
	case "group-id":
		t.GroupID = uint(temp)
	}

	return errors.New("There is no to-one relationship with the name " + name)
}
