package models

import (
	"errors"
	"payup/app-error"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/manyminds/api2go/jsonapi"
)

// Transaction model
type Transaction struct {
	ID                uint `json:"id"`
	Type              string
	Amount            int
	Memo              string
	RecipientID       uint   `jsonapi:"name=recipientId"`
	SenderID          uint   `jsonapi:"name=senderId"`
	RelatedObjectType string `jsonapi:"name=relatedObjectType"`
	RelatedObjectID   uint   `jsonapi:"name=relatedObjectId"`
	CreatorID         uint   `jsonapi:"name=creatorId"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `jsonapi:"-"`

	Recipient User `jsonapi:"-" sql:"-"`
	Sender    User `jsonapi:"-" sql:"-"`
	Creator   User `jsonapi:"-" sql:"-"`
}

// BeforeUpdate ensures that friendship balance is kept in sync
func (t *Transaction) BeforeUpdate(db *gorm.DB) (err error) {
	if t.RelatedObjectType != "Friendship" {
		return
	}

	var curTransaction Transaction

	db.First(&curTransaction, t.ID)

	ReverseTransaction(&curTransaction, db)
	// Now the AfterSave callback will use the new updated transaction
	// and update the balance accordingly
	return
}

// AfterSave increments balance on FriendshipData
func (t *Transaction) AfterSave(db *gorm.DB) (err error) {
	if t.RelatedObjectType != "Friendship" {
		return
	}
	// Transaction is related to a Friendship
	var fd FriendshipData
	db.First(&fd, t.RelatedObjectID)

	switch {
	case fd.PositiveUserID == t.RecipientID:
		fd.Balance -= t.Amount
	case fd.PositiveUserID == t.SenderID:
		fd.Balance += t.Amount
	}

	db.Save(&fd)

	return
}

// AfterDelete ensures that friendship balance is reversed (as if this transaction never occurred)
func (t *Transaction) AfterDelete(db *gorm.DB) (err error) {
	if t.RelatedObjectType != "Friendship" {
		return
	}
	ReverseTransaction(t, db)
	return
}

// ReverseTransaction this function will take a transaction amount and Type
// and users to reverse the transaction on the balance
func ReverseTransaction(t *Transaction, db *gorm.DB) {
	// Transaction is related to a Friendship
	var fd FriendshipData
	db.First(&fd, t.RelatedObjectID)

	// Reverse the old transaction
	switch {
	case fd.PositiveUserID == t.RecipientID:
		fd.Balance += t.Amount
	case fd.PositiveUserID == t.SenderID:
		fd.Balance -= t.Amount
	}

	// Save the new FriendshipData
	db.Save(&fd)
}

// Validate the transaction and return a boolean and appError
func (t Transaction) Validate() (bool, appError.Err) {
	if t.Type != "Bill" && t.Type != "Payment" {
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

	if t.SenderID == 0 {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The transaction senderId cannot be 0 or empty"
		return false, invalidID
	}

	if t.RecipientID == 0 {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The transaction recipientId cannot be 0 or empty"
		return false, invalidID
	}

	if t.RelatedObjectID == 0 {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The transaction relatedObjectID cannot be 0 or empty"
		return false, invalidID
	}

	if t.RelatedObjectType != "Group" && t.RelatedObjectType != "Friendship" {
		invalidType := appError.InvalidParams
		invalidType.Detail = "The transaction must have a valid relatedObjectType"
		return false, invalidType
	}

	return true, appError.Err{}
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

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
			Name: "recipient",
		},
		{
			Type: "users",
			Name: "sender",
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
		ID:   strconv.FormatUint(uint64(t.RecipientID), 10),
		Type: "users",
		Name: "recipient",
	})

	result = append(result, jsonapi.ReferenceID{
		ID:   strconv.FormatUint(uint64(t.SenderID), 10),
		Type: "users",
		Name: "sender",
	})
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (t Transaction) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	result = append(result, t.Recipient)
	result = append(result, t.Sender)
	result = append(result, t.Creator)

	return result
}

// SetToOneReferenceID sets the reference ID and satisfies the jsonapi.UnmarshalToOneRelations interface
func (t *Transaction) SetToOneReferenceID(name, ID string) error {
	temp, err := strconv.ParseUint(ID, 10, 64)

	if err != nil {
		return err
	}

	switch name {
	case "recipient":
		t.RecipientID = uint(temp)
	case "sender":
		t.SenderID = uint(temp)
	case "creator":
		t.CreatorID = uint(temp)
	}

	return errors.New("There is no to-one relationship with the name " + name)
}
