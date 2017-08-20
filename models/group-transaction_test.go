package models

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tbbr/tbbr-api/app-error"
)

type GroupTransactionModelTestSuite struct {
	suite.Suite
}

func TestGroupTransactionModelTestSuite(t *testing.T) {
	suite.Run(t, new(GroupTransactionModelTestSuite))
}

func (s *TransactionModelTestSuite) TestgtGenerateSplits_SimpleOnePart() {
	splits := gtGenerateSimpleSplits(1000, 1)
	assert.Equal(s.T(), 1, len(splits))
	assert.Equal(s.T(), int64(1000), splits[0])
}

func (s *TransactionModelTestSuite) TestgtGenerateSplits_SimpleParts() {
	splits := gtGenerateSimpleSplits(2000, 9)
	assert.Equal(s.T(), pq.Int64Array{223, 223, 222, 222, 222, 222, 222, 222, 222}, splits)
}

func (s *TransactionModelTestSuite) TestgtGenerateSplits_SplitCountZero() {
	splits := gtGenerateSimpleSplits(2000, 0)
	assert.Equal(s.T(), pq.Int64Array(nil), splits)
}

func (s *TransactionModelTestSuite) TestgtGenerateSplitAmounts_SimpleOnePart() {
	splits := gtGenerateSplitAmounts(2000, pq.Int64Array{3})
	assert.Equal(s.T(), 1, len(splits))
	assert.Equal(s.T(), int64(2000), splits[0])
}

func (s *TransactionModelTestSuite) TestgtGenerateSplitAmounts_MultiPart() {
	splits := gtGenerateSplitAmounts(2000, pq.Int64Array{3, 5, 1})
	assert.Equal(s.T(), 3, len(splits))
	assert.Equal(s.T(), int64(667), splits[0])
	assert.Equal(s.T(), int64(1111), splits[1])
	assert.Equal(s.T(), int64(222), splits[2])
}

func (s *TransactionModelTestSuite) TestGetRecipientSplitAmounts_SimpleSplit() {
	gt := GroupTransaction{
		Amount:             1501,
		Memo:               "test",
		SenderIDs:          pq.Int64Array{1, 2},
		RecipientIDs:       pq.Int64Array{5, 6, 7},
		SenderSplits:       pq.Int64Array{1400, 101},
		RecipientSplits:    pq.Int64Array{1, 3, 2},
		RecipientSplitType: "splitPart",
		SenderSplitType:    "normal",
		GroupID:            2,
		CreatorID:          4,
	}

	isValid, err := gt.Validate()
	assert.True(s.T(), isValid)
	assert.Equal(s.T(), appError.Err{}, err)

	recipientSplits := gt.GetRecipientSplitAmounts()
	assert.Equal(s.T(), pq.Int64Array{251, 750, 500}, recipientSplits)

	senderSplits := gt.GetSenderSplitAmounts()
	assert.Equal(s.T(), pq.Int64Array{1400, 101}, senderSplits)
}

func (s *TransactionModelTestSuite) TestGetSenderSplitAmounts_SimpleSplit() {
	gt := GroupTransaction{
		Amount:             20001,
		Memo:               "test",
		SenderIDs:          pq.Int64Array{1, 2},
		RecipientIDs:       pq.Int64Array{5, 6, 7},
		SenderSplits:       pq.Int64Array{5, 15},
		RecipientSplits:    pq.Int64Array{10000, 1, 10000},
		RecipientSplitType: "normal",
		SenderSplitType:    "splitPart",
		GroupID:            2,
		CreatorID:          4,
	}

	isValid, err := gt.Validate()
	assert.True(s.T(), isValid)
	assert.Equal(s.T(), appError.Err{}, err)

	senderSplits := gt.GetSenderSplitAmounts()
	assert.Equal(s.T(), pq.Int64Array{5001, 15000}, senderSplits)

	recipientSplits := gt.GetRecipientSplitAmounts()
	assert.Equal(s.T(), pq.Int64Array{10000, 1, 10000}, recipientSplits)
}
