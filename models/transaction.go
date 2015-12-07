package models

import (
	"errors"
	"payup/app-error"
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
	DeletedAt     *time.Time `jsonapi:"-"`

	Creator     User `jsonapi:"-" sql:"-"`
	RelatedUser User `jsonapi:"-" sql:"-"`
}

// Validate the transaction and return a boolean and appError
func (t Transaction) Validate() (bool, appError.Err) {
	if t.Type != "Borrow" && t.Type != "Lend" {
		invalidType := appError.InvalidParams
		invalidType.Detail = "The transaction type is invalid"
		return false, invalidType
	}

	// Maximum amount of $10,000
	if t.Amount > 1000000 || t.Amount < 0 {
		invalidAmount := appError.InvalidParams
		invalidAmount.Detail = "The transaction amount is out of range"
		return false, invalidAmount
	}

	if len([]rune(t.Memo)) > 140 {
		invalidMemo := appError.InvalidParams
		invalidMemo.Detail = "The transaction memo must be less than or equal to 140 characters"
		return false, invalidMemo
	}

	if t.RelatedUserID == 0 {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The transaction relatedUserId cannot be 0"
		return false, invalidID
	}

	return true, appError.Err{}
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

// GetReferences returns all related structs to transactions
func (t Transaction) GetReferences() []jsonapi.Reference {
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
