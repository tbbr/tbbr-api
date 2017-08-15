package repositories

import (
	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"
)

// GroupMemberRepository struct just there for namespacing
type GroupMemberRepository struct {
}

// NewGroupMemberRepository returns a GroupRepository
func NewGroupMemberRepository() *GroupMemberRepository {
	var newGroupMemberRepository GroupMemberRepository
	return &newGroupMemberRepository
}

// Create takes a GroupMember struct and persists it into DB
func (r *GroupMemberRepository) Create(gm models.GroupMember) (*models.GroupMember, error) {
	if dbc := database.DBCon.Create(&gm); dbc.Error != nil {
		return nil, dbc.Error
	}
	return &gm, nil
}
