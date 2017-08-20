package models

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/tbbr/tbbr-api/app-error"
	"github.com/tbbr/tbbr-api/database"
)

// GroupTransaction is similar to a Transaction model except that
// the transaction can have multiple senders and multiple recipients
type GroupTransaction struct {
	ID                 uint          `json:"-"`
	Amount             uint          `json:"amount"`
	Memo               string        `json:"memo"`
	SenderIDs          pq.Int64Array `gorm:"type:integer[]" json:"senderIds"`
	RecipientIDs       pq.Int64Array `gorm:"type:integer[]" json:"recipientIds"`
	SenderSplits       pq.Int64Array `gorm:"type:integer[]" json:"senderSplits"`
	RecipientSplits    pq.Int64Array `gorm:"type:integer[]" json:"recipientSplits"`
	SenderSplitType    string        `json:"senderSplitType"`
	RecipientSplitType string        `json:"recipientSplitType"`
	GroupID            uint          `json:"groupId"`
	CreatorID          uint          `json:"-"`
	CreatedAt          time.Time     `json:"createdAt"`
	UpdatedAt          time.Time     `json:"updatedAt"`
	DeletedAt          *time.Time    `json:"-"`

	Senders    []User `json:"-" sql:"-"`
	Recipients []User `json:"-" sql:"-"`
	Group      Group  `json:"-" sql:"-"`
	Creator    User   `json:"-" sql:"-"`
}

func gtGenerateSimpleSplits(amount uint, splitCount int64) pq.Int64Array {
	amountLeft := int64(amount)
	var splits pq.Int64Array
	if splitCount == 0 {
		return nil
	}
	for splitCount != 0 {
		res := math.Ceil(float64(amountLeft) / float64(splitCount))
		amountLeft -= int64(res)
		splitCount--
		splits = append(splits, int64(res))
	}
	return splits
}

func gtGenerateSplitAmounts(amount uint, splitParts pq.Int64Array) pq.Int64Array {
	splitPartTotal := int64(0)
	for i := range splitParts {
		splitPartTotal += splitParts[i]
	}
	simpleSplits := gtGenerateSimpleSplits(amount, splitPartTotal)
	consolidatedSplits := make(pq.Int64Array, len(splitParts))
	cur := int64(0)

	maxSplitPart := int64(0)
	for i := range splitParts {
		if splitParts[i] > maxSplitPart {
			maxSplitPart = splitParts[i]
		}
	}
	// Keep running until we've consumed all split parts
	for cur < splitPartTotal {
		for i := 0; i < len(consolidatedSplits); i++ {
			if splitParts[i] > 0 {
				consolidatedSplits[i] += simpleSplits[cur]
				splitParts[i]--
				cur++
			}
			// if we've consumed all the split parts, then exit
			if cur == splitPartTotal {
				break
			}
		}
	}
	return consolidatedSplits
}

func (gt *GroupTransaction) GetSenderSplitAmounts() pq.Int64Array {
	if gt.SenderSplitType == "normal" {
		return gt.SenderSplits
	} else if gt.SenderSplitType == "splitPart" {
		return gtGenerateSplitAmounts(gt.Amount, gt.SenderSplits)
	}
	return nil
}

func (gt *GroupTransaction) GetRecipientSplitAmounts() pq.Int64Array {
	if gt.RecipientSplitType == "normal" {
		return gt.RecipientSplits
	} else if gt.RecipientSplitType == "splitPart" {
		return gtGenerateSplitAmounts(gt.Amount, gt.RecipientSplits)
	}
	return nil
}

// BeforeUpdate ensures that friendship balance is kept in sync
func (gt *GroupTransaction) BeforeUpdate(db *gorm.DB) (err error) {
	var cur GroupTransaction
	db.First(&cur, gt.ID)
	ReverseGroupTransaction(&cur, db)
	// Now the AfterSave callback will use the new updated transaction
	// and update the balance accordingly
	return
}

// AfterSave increments balance on FriendshipData
func (gt *GroupTransaction) AfterSave(db *gorm.DB) (err error) {
	// Transaction is related to a Friendship
	var senderMembers []GroupMember
	var recipientMembers []GroupMember

	database.DBCon.Where("user_id in (?) AND group_id = ?", gt.SenderIDs, gt.GroupID).Find(&senderMembers)
	database.DBCon.Where("user_id in (?) AND group_id = ?", gt.RecipientIDs, gt.GroupID).Find(&recipientMembers)

	senderSplitAmounts := gt.GetSenderSplitAmounts()
	recipientSplitAmounts := gt.GetRecipientSplitAmounts()

	for i := range senderMembers {
		senderMembers[i].AmountSent += uint(senderSplitAmounts[i])
		database.DBCon.Model(&senderMembers[i]).Update("amount_sent", senderMembers[i].AmountSent)
	}

	for i := range recipientMembers {
		recipientMembers[i].AmountReceived += uint(recipientSplitAmounts[i])
		database.DBCon.Model(&recipientMembers[i]).Update("amount_received", recipientMembers[i].AmountReceived)
	}
	// gt.sendNotification()
	return nil
}

// AfterDelete ensures that friendship balance is reversed (as if this transaction never occurred)
func (gt *GroupTransaction) AfterDelete(db *gorm.DB) (err error) {
	ReverseGroupTransaction(gt, db)
	return nil
}

//
// ReverseGroupTransaction will take a GroupTransaction amount
// and reverse it
func ReverseGroupTransaction(gt *GroupTransaction, db *gorm.DB) {
	var senderMembers []GroupMember
	var recipientMembers []GroupMember

	database.DBCon.Where("user_id in (?) AND group_id = ?", gt.SenderIDs, gt.GroupID).Find(&senderMembers)
	database.DBCon.Where("user_id in (?) AND group_id = ?", gt.RecipientIDs, gt.GroupID).Find(&recipientMembers)

	senderSplitAmounts := gt.GetSenderSplitAmounts()
	recipientSplitAmounts := gt.GetRecipientSplitAmounts()

	for i := range senderMembers {
		senderMembers[i].AmountSent -= uint(senderSplitAmounts[i])
		database.DBCon.Model(senderMembers[i]).Update("amount_sent", senderMembers[i].AmountSent)
	}

	for i := range recipientMembers {
		recipientMembers[i].AmountReceived -= uint(recipientSplitAmounts[i])
		database.DBCon.Model(recipientMembers[i]).Update("amount_received", recipientMembers[i].AmountReceived)
	}
}

// Validate the transaction and return a boolean and appError
func (gt GroupTransaction) Validate() (bool, appError.Err) {
	// Maximum amount of $100,000
	if gt.Amount > 10000000 || gt.Amount < 0 {
		invalidAmount := appError.InvalidParams
		invalidAmount.Detail = "The groupTransaction amount must be between $0 and $100,000"
		return false, invalidAmount
	}

	if len([]rune(gt.Memo)) > 140 {
		invalidMemo := appError.InvalidParams
		invalidMemo.Detail = "The groupTransaction memo must be less than or equal to 140 characters"
		return false, invalidMemo
	}

	if gt.GroupID == 0 {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The groupTransaction groupId must not be empty or 0"
		return false, invalidID
	}

	if len(gt.SenderIDs) == 0 {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The groupTransaction must have at least one sender"
		return false, invalidID
	}

	if len(gt.RecipientIDs) == 0 {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The groupTransaction must have at least one recipient"
		return false, invalidID
	}

	if len(gt.RecipientSplits) != len(gt.RecipientIDs) {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The groupTransaction len recipientSplits must match len of recipientIds"
		return false, invalidID
	}

	if len(gt.SenderSplits) != len(gt.SenderIDs) {
		invalidID := appError.InvalidParams
		invalidID.Detail = "The groupTransaction len senderSplits must match len of senderIds"
		return false, invalidID
	}

	if gt.SenderSplitType != "splitPart" && gt.SenderSplitType != "normal" {
		invalidType := appError.InvalidParams
		invalidType.Detail = "The groupTransaction senderSplitType is invalid, must be one of (splitPart, normal)"
		return false, invalidType
	}

	if gt.RecipientSplitType != "splitPart" && gt.RecipientSplitType != "normal" {
		invalidType := appError.InvalidParams
		invalidType.Detail = "The groupTransaction recipientSplitType is invalid, must be one of (splitPart, normal)"
		return false, invalidType
	}

	return true, appError.Err{}
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

// GetID returns a stringified version of an ID
func (gt GroupTransaction) GetID() string {
	return strconv.FormatUint(uint64(gt.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (gt *GroupTransaction) SetID(id string) error {
	groupTransactionID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	gt.ID = uint(groupTransactionID)
	return nil
}

// GetReferences returns all related structs to groupTransactions
func (gt GroupTransaction) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "creator",
		},
		{
			Type: "users",
			Name: "senders",
		},
		{
			Type: "users",
			Name: "recipients",
		},
	}
}

// GetReferencedIDs satisfies the jsonapi.MarshalLinkedRelations interface
func (gt GroupTransaction) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	result = append(result, jsonapi.ReferenceID{
		ID:   strconv.FormatUint(uint64(gt.CreatorID), 10),
		Type: "users",
		Name: "creator",
	})

	for _, user := range gt.Senders {
		result = append(result, jsonapi.ReferenceID{
			ID:   user.GetID(),
			Type: "users",
			Name: "senders",
		})
	}

	for _, user := range gt.Recipients {
		result = append(result, jsonapi.ReferenceID{
			ID:   user.GetID(),
			Type: "users",
			Name: "recipients",
		})
	}
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (gt GroupTransaction) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	result = append(result, gt.Creator)
	for key := range gt.Senders {
		result = append(result, gt.Senders[key])
	}
	for key := range gt.Recipients {
		result = append(result, gt.Recipients[key])
	}
	return result
}

// SetToOneReferenceID sets the reference ID and satisfies the jsonapi.UnmarshalToOneRelations interface
func (gt *GroupTransaction) SetToOneReferenceID(name, ID string) error {
	temp, err := strconv.ParseUint(ID, 10, 64)

	if err != nil {
		return err
	}

	switch name {
	case "creator":
		gt.CreatorID = uint(temp)
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}
