package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/speps/go-hashids"
)

// Group model that users wil use
type Group struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Users        []User        `gorm:"many2many:group_users;" json:"users"`
	Transactions []Transaction `json:"transactions"`
	HashID       string        `json:"hashId"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	DeletedAt    *time.Time    `json:"deletedAt"`
}

// AfterCreate generates a HashID for a Group based on it's numeric ID field
func (g *Group) AfterCreate(db *gorm.DB) (err error) {
	hd := hashids.NewData()
	hd.Salt = "9398dfajsie288sawiehg"
	hd.MinLength = 6
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
