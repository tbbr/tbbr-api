package models

import (
	"time"
)

// GroupTransaction is similar to a Transaction model except that
// the transaction can have multiple senders and multiple recipients
type GroupTransaction struct {
	ID                 uint       `json:"-"`
	Amount             string     `json:"name"`
	Senders            []uint     `json:"senders"`
	Recipients         []uint     `json:"recipients"`
	SenderSplits       []uint     `json:"senderSplits"`
	RecipientSplits    []uint     `json:"recipientSplits"`
	SenderSplitType    string     `json:"senderSplitType"`
	RecipientSplitType string     `json:"recipientSplitType"`
	GroupID            string     `json:"-"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `json:"-"`

	SendersContent    []User  `json:"-" sql:"-"`
	RecipientsContent []User  `json:"-" sql:"-"`
	GroupContent      []Group `json:"-" sql:"-"`
}
