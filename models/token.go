package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// Token model
type Token struct {
	ID                uint   `json:"id"`
	Category          string `json:"category"`
	AccessToken       string `json:"accessToken"`
	RefreshToken      string `json:"refreshToken"`
	RefreshExpiration time.Time
	AuthExpiration    time.Time
	Expired           bool       `json:"expired"`
	UserID            uint       `json:"userId"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt"`
}

// GetID returns a stringified version of an ID
func (t Token) GetID() string {
	return strconv.FormatUint(uint64(t.ID), 10)
}

// BeforeCreate generates a access and refresh tokens
// and expiry dates
func (t *Token) BeforeCreate(db *gorm.DB) (err error) {
	t.AccessToken = uuid.NewV4().String()
	t.RefreshToken = uuid.NewV4().String()
	t.RefreshExpiration = time.Now().AddDate(0, 0, 3)
	t.AuthExpiration = time.Now().AddDate(0, 0, 1)
	t.Expired = false
	return
}
