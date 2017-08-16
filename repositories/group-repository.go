package repositories

import (
	"github.com/tbbr/tbbr-api/app-error"
	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"
)

// GroupRepository struct just there for namespacing
type GroupRepository struct {
}

// NewGroupRepository returns a GroupRepository
func NewGroupRepository() *GroupRepository {
	var newGroupRepository GroupRepository
	return &newGroupRepository
}

// List is a function that returns a list of groups
// @params
//  user: Get the groups of this user
//  limit: return x amount of groups
//  offset: skip the first x groups
// @returns
//  groups: a list of groups
func (r *GroupRepository) List(user models.User, limit int, offset int) []models.Group {
	var groups []models.Group
	database.DBCon.Model(&user).Preload("Users").Related(&groups, "Groups")
	return groups
}

// Get is a function that returns a specific group
// @params
//  groupID: the id of the group returned
// @returns
//  group: a group that has groupID as it's own id
func (r *GroupRepository) Get(groupID uint) (*models.Group, *appError.Err) {
	var group models.Group
	if database.DBCon.First(&group, groupID).RecordNotFound() {
		return nil, &appError.RecordNotFound
	}

	database.DBCon.Model(&group).Related(&group.GroupMembers)
	return &group, nil
}

// Create takes a group struct and persists it into the DB
func (r *GroupRepository) Create(g models.Group) (*models.Group, error) {
	if dbc := database.DBCon.Create(&g); dbc.Error != nil {
		return nil, dbc.Error
	}
	return &g, nil
}

// Update takes a group struct and updates specific fields
func (r *GroupRepository) Update(g models.Group) (*models.Group, error) {
	dbc := database.DBCon.Model(&g).Update(map[string]interface{}{
		"name":        g.Name,
		"description": g.Description,
	})
	if dbc.Error != nil {
		return nil, dbc.Error
	}
	return &g, nil
}

// AddGroupMember takes a groupID and a userID and creates a groupMember
func (r *GroupRepository) AddGroupMember(groupID uint, userID uint) error {
	gmr := NewGroupMemberRepository()
	var gm models.GroupMember

	gm.SetDefault()
	gm.GroupID = groupID
	gm.UserID = userID

	_, err := gmr.Create(gm)
	if err != nil {
		return err
	}
	return nil
}
