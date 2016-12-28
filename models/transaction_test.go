package models

import (
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tbbr/tbbr-api/app-error"
)

type TransactionModelTestSuite struct {
	suite.Suite
}

func TestTransactionModelTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionModelTestSuite))
}

func (s *TransactionModelTestSuite) TestUnmarshal_SimpleFields() {
	jsonByte := []byte(`{
        "data": {
            "attributes": {
                "type": "Bill",
                "amount": 4200,
                "isSettled": false,
                "memo": "test",
                "status": "Confirmed",
                "relatedObjectType": "Friendship",
                "relatedObjectId": 1
            },
            "type": "transactions"
        }
    }`)
	var t Transaction
	err := jsonapi.Unmarshal(jsonByte, &t)
	assert.Nil(s.T(), err, "JSON unmarshalling error must be nil")

	assert.Empty(s.T(), t.RecipientID, "RecipientID should not be set")
	assert.Empty(s.T(), t.SenderID, "SenderID should not be set")
	assert.Empty(s.T(), t.CreatorID, "CreatorID should not be set")

	assert.Equal(s.T(), t.Type, "Bill", "Type should be Bill")
	assert.Equal(s.T(), t.Amount, 4200, "Amount should be 4200")
	assert.Equal(s.T(), t.IsSettled, false, "IsSettled should be false")
	assert.Equal(s.T(), t.Memo, "test", "Memo should be test")
	assert.Equal(s.T(), t.Status, "Confirmed", "Status should be Confirmed")
	assert.Equal(s.T(), t.RelatedObjectID, uint(1), "RelatedObjectID is 1")
	assert.Equal(s.T(), t.RelatedObjectType, "Friendship", "RelatedObjectType is Friendship")

	isValid, appErr := t.Validate()
	assert.False(s.T(), isValid)
	assert.Equal(s.T(), appErr.Code, "2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams")
}

func (s *TransactionModelTestSuite) TestUnmarshal_Relations() {
	jsonByte := []byte(`{
        "data": {
            "attributes": {
                "type": "Bill",
                "amount": 4200,
                "isSettled": false,
                "memo": "test",
                "status": "Confirmed",
                "relatedObjectType": "Friendship",
                "relatedObjectId": 1
            },
            "relationships": {
                "recipient": {
                    "data": {
                        "id": "2",
                        "type": "users"
                    }
                },
                "sender": {
                    "data": {
                        "id": "1",
                        "type": "users"
                    }
                }
            },
            "type": "transactions"
        }
    }`)

	var t Transaction
	err := jsonapi.Unmarshal(jsonByte, &t)
	assert.Nil(s.T(), err, "JSON unmarshalling error must be nil")
	assert.Equal(s.T(), t.RecipientID, uint(2), "RecipientID should be set to 2")
	assert.Equal(s.T(), t.SenderID, uint(1), "SenderID should be set to 1")
	assert.Empty(s.T(), t.CreatorID, "CreatorID should not be set")

	isValid, appErr := t.Validate()
	assert.True(s.T(), isValid)
	assert.Equal(s.T(), appErr, (appError.Err{}), "the error must be zeroed out")
}

func (s *TransactionModelTestSuite) TestGetFormattedAmount_Normal() {
	t := Transaction{
		ID:     1,
		Amount: 5210,
	}

	t.GetFormattedAmount()
	assert.Equal(s.T(), t.GetFormattedAmount(), "$52.10", "GetFormattedAmount returns correct currency")
}

func (s *TransactionModelTestSuite) TestGetFormattedAmount_NotSpecified() {
	t := Transaction{
		ID: 1,
	}

	t.GetFormattedAmount()
	assert.Equal(s.T(), t.GetFormattedAmount(), "$0.00", "GetFormattedAmount returns correct currency")
}

func (s *TransactionModelTestSuite) TestValidate_InvalidType() {
	t := Transaction{
		ID:                1,
		Type:              "ABCD",
		Status:            "Confirmed",
		Amount:            8910,
		Memo:              "This is a simple test",
		IsSettled:         false,
		RecipientID:       uint(1),
		SenderID:          uint(2),
		RelatedObjectType: "Friendship",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction type is invalid", "Detail should be about transaction type!")
}

func (s *TransactionModelTestSuite) TestValidate_InvalidStatus() {
	t := Transaction{
		ID:                1,
		Type:              "Payback",
		Status:            "NOTVALID",
		Amount:            8910,
		Memo:              "This is a simple test",
		IsSettled:         false,
		RecipientID:       uint(1),
		SenderID:          uint(2),
		RelatedObjectType: "Friendship",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction status is invalid", "Detail should be about transaction status!")
}

func (s *TransactionModelTestSuite) TestValidate_InvalidAmountTooLarge() {
	t := Transaction{
		ID:                1,
		Type:              "Bill",
		Status:            "Rejected",
		Amount:            10000001, // $100,000.01
		Memo:              "This is a simple test",
		IsSettled:         false,
		RecipientID:       uint(1),
		SenderID:          uint(2),
		RelatedObjectType: "Friendship",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction amount must be between $0 and $100,000", "Detail should be about Amount")
}

func (s *TransactionModelTestSuite) TestValidate_InvalidAmountTooSmall() {
	t := Transaction{
		ID:                1,
		Type:              "Bill",
		Status:            "Rejected",
		Amount:            -20,
		Memo:              "This is a simple test",
		IsSettled:         false,
		RecipientID:       uint(1),
		SenderID:          uint(2),
		RelatedObjectType: "Friendship",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction amount must be between $0 and $100,000", "Detail should be about Amount")
}

func (s *TransactionModelTestSuite) TestValidate_InvalidMemo() {
	t := Transaction{
		ID:                1,
		Type:              "Payback",
		Status:            "Pending",
		Amount:            8910,
		Memo:              "This is a memo that will be longer than 140 characters, since that's the limit we should get an invalid parameters transaction. abcdefghijklm", // exactly 141 characters long
		IsSettled:         false,
		RecipientID:       uint(1),
		SenderID:          uint(2),
		RelatedObjectType: "Friendship",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction memo must be less than or equal to 140 characters", "Detail should be about transaction memo!")
}

func (s *TransactionModelTestSuite) TestValidate_EmptySenderID() {
	t := Transaction{
		ID:                1,
		Type:              "Payback",
		Status:            "Pending",
		Amount:            8910,
		Memo:              "This is a memo that will be longer than 140 characters, since that's the limit we should get an invalid parameters transaction. abcdefghijkl", // exactly 140 characters long
		IsSettled:         false,
		RecipientID:       uint(1),
		RelatedObjectType: "Friendship",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction senderId cannot be 0 or empty", "Detail should be about transaction senderId!")
}

func (s *TransactionModelTestSuite) TestValidate_EmptyRecipientID() {
	t := Transaction{
		ID:                1,
		Type:              "Payback",
		Status:            "Pending",
		Amount:            8910,
		Memo:              "This is a memo that will be longer than 140 characters, since that's the limit we should get an invalid parameters transaction. abcdefghijkl", // exactly 140 characters long
		IsSettled:         false,
		SenderID:          uint(2),
		RelatedObjectType: "Group",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction recipientId cannot be 0 or empty", "Detail should be about transaction recipientId!")
}

func (s *TransactionModelTestSuite) TestValidate_SenderAndRecipientIsSame() {
	t := Transaction{
		ID:                1,
		Type:              "Payback",
		Status:            "Pending",
		Amount:            8910,
		Memo:              "This is a memo that will be longer than 140 characters, since that's the limit we should get an invalid parameters transaction. abcdefghijkl", // exactly 140 characters long
		IsSettled:         false,
		SenderID:          uint(2),
		RecipientID:       uint(2),
		RelatedObjectType: "Friendship",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction recipient and sender cannot be the same", "Detail should be about transaction recipientId and senderId")
}

func (s *TransactionModelTestSuite) TestValidate_EmptyRelatedObjectID() {
	t := Transaction{
		ID:                1,
		Type:              "Payback",
		Status:            "Pending",
		Amount:            8910,
		Memo:              "This is a memo that will be longer than 140 characters, since that's the limit we should get an invalid parameters transaction. abcdefghijkl", // exactly 140 characters long
		IsSettled:         false,
		SenderID:          uint(2),
		RecipientID:       uint(1),
		RelatedObjectType: "Friendship",
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction relatedObjectID cannot be 0 or empty", "Detail should be about transaction relatedObjectID")
}

func (s *TransactionModelTestSuite) TestValidate_InvalidRelatedObjectType() {
	t := Transaction{
		ID:                1,
		Type:              "Payback",
		Status:            "Pending",
		Amount:            8910,
		Memo:              "This is a memo that will be longer than 140 characters, since that's the limit we should get an invalid parameters transaction. abcdefghijkl", // exactly 140 characters long
		IsSettled:         false,
		SenderID:          uint(2),
		RecipientID:       uint(1),
		RelatedObjectType: "ABCD",
		RelatedObjectID:   uint(2),
	}

	isValid, appErr := t.Validate()

	assert.False(s.T(), isValid, "transaction should not be valid")
	assert.Equal(s.T(), appErr.Code, "2000", "Code should be 2000")
	assert.Equal(s.T(), appErr.Title, "InvalidParams", "Title should be InvalidParams")
	assert.Equal(s.T(), appErr.Detail, "The transaction must have a valid relatedObjectType", "Detail should be about transaction relatedObjectType")
}

func (s *TransactionModelTestSuite) TestValidate_ValidTransaction() {
	t := Transaction{
		ID:                1,
		Type:              "Payback",
		Status:            "Pending",
		Amount:            8910,
		Memo:              "This is a memo that will be longer than 140 characters, since that's the limit we should get an invalid parameters transaction. abcdefghijkl", // exactly 140 characters long
		IsSettled:         false,
		SenderID:          uint(2),
		RecipientID:       uint(1),
		RelatedObjectType: "Group",
		RelatedObjectID:   uint(5),
	}

	isValid, appErr := t.Validate()

	assert.True(s.T(), isValid, "transaction should be valid")
	assert.Equal(s.T(), appErr, (appError.Err{}), "appErr should be zeroed (no error)")
}
