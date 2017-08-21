package repositories

import (
	"github.com/tbbr/tbbr-api/app-error"
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

// List is a function that returns a list of groupMembers
// @params
//  groupID: members for x group
//  limit: return x amount of members
//  offset: skip the first x members
// @returns
//  groupMembers: a list of groupMembers
func (r *GroupMemberRepository) List(groupID uint, limit int, offset int) []models.GroupMember {
	var groupMembers []models.GroupMember

	if groupID <= 0 {
		return nil
	}

	database.DBCon.
		Where("group_id = ?", groupID).
		Preload("User").
		Find(&groupMembers)

	return groupMembers
}

// Get is a function that returns a specific groupMember
// @params
//  id: the id of the member returned
// @returns
//  groupMember
func (r *GroupMemberRepository) Get(id uint) (*models.GroupMember, *appError.Err) {
	var gm models.GroupMember
	if database.DBCon.First(&gm, id).RecordNotFound() {
		return nil, &appError.RecordNotFound
	}

	database.DBCon.First(&gm.User, gm.UserID)
	return &gm, nil
}

// Create takes a GroupMember struct and persists it into DB
func (r *GroupMemberRepository) Create(gm models.GroupMember) (*models.GroupMember, *appError.Err) {
	if dbc := database.DBCon.Create(&gm); dbc.Error != nil {
		dbError := appError.DatabaseError
		dbError.Detail = dbc.Error.Error()
		return nil, &dbError
	}
	return &gm, nil
}
