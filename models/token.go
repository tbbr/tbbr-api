package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// Token model
type Token struct {
	ID                uint       `json:"id"`
	Category          string     `json:"category"`
	AccessToken       string     `json:"accessToken"`
	RefreshToken      string     `json:"refreshToken"`
	RefreshExpiration time.Time  `json:"refreshExpiration"`
	AuthExpiration    time.Time  `json:"authExpiration"`
	UserID            uint       `json:"userId"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         *time.Time `json:"deletedAt"`
}

// BeforeCreate generates access and refresh tokens
// and expiry dates
func (t *Token) BeforeCreate(db *gorm.DB) (err error) {
	t.AccessToken = uuid.NewV4().String()
	t.RefreshToken = uuid.NewV4().String()
	t.RefreshExpiration = time.Now().AddDate(0, 0, 3) // 3 days from now
	t.AuthExpiration = time.Now().AddDate(0, 0, 1)    // 1 day from now
	return
}

// Expired function returns true if the AccessToken has expired, and false otherwise
func (t Token) Expired() bool {
	return time.Now().After(t.AuthExpiration)
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

// GetID returns a stringified version of an ID
func (t Token) GetID() string {
	return strconv.FormatUint(uint64(t.ID), 10)
}
