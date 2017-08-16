package repositories

import (
	"github.com/tbbr/tbbr-api/app-error"
	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"
)

// GroupRepository struct just there for namespacing
type GroupTransactionRepository struct {
}

// NewGroupRepository returns a GroupRepository
func NewGroupTransactionRepository() *GroupTransactionRepository {
	var newGroupTransactionRepository GroupTransactionRepository
	return &newGroupTransactionRepository
}

// List is a function that returns a list of groupTransactions
// @params
//  groupID: transactions for x group
//  limit: return x amount of groups
//  offset: skip the first x groups
// @returns
//  groups: a list of groups
func (r *GroupTransactionRepository) List(groupID uint, limit int, offset int) []models.GroupTransaction {
	var groupTransactions []models.GroupTransaction

	if groupID <= 0 {
		return nil
	}

	if groupID > 0 {
		database.DBCon.
			Where("group_id = ?", groupID).
			Order("created_at desc").
			Find(&groupTransactions)
	}

	// Get related data for transactions
	// TODO: n + 1 query problem here, should use eager loading in someway
	for i := range groupTransactions {
		database.DBCon.Where("id in (?)", groupTransactions[i].SenderIDs).Find(&groupTransactions[i].Senders)
		database.DBCon.Where("id in (?)", groupTransactions[i].RecipientIDs).Find(&groupTransactions[i].Recipients)
		database.DBCon.First(&groupTransactions[i].Creator, groupTransactions[i].CreatorID)
	}

	return groupTransactions
}

// Get is a function that returns a specific groupTransaction
// @params
//  groupTransactionID: the id of the group returned
// @returns
//  group: a group that has groupID as it's own id
func (r *GroupTransactionRepository) Get(groupTransactionID uint) (*models.GroupTransaction, *appError.Err) {
	var gt models.GroupTransaction
	if database.DBCon.First(&gt, groupTransactionID).RecordNotFound() {
		return nil, &appError.RecordNotFound
	}

	database.DBCon.Where("id in (?)", gt.SenderIDs).Find(&gt.Senders)
	database.DBCon.Where("id in (?)", gt.RecipientIDs).Find(&gt.Recipients)
	database.DBCon.First(&gt.Creator, gt.CreatorID)
	return &gt, nil
}

// Create takes a groupTransaction struct and persists it into the DB
func (r *GroupTransactionRepository) Create(gt models.GroupTransaction) (*models.GroupTransaction, *appError.Err) {
	// Validate our new transaction
	isValid, errApp := gt.Validate()

	if isValid == false {
		return nil, &errApp
	}

	if dbc := database.DBCon.Create(&gt); dbc.Error != nil {
		dbError := appError.DatabaseError
		dbError.Detail = dbc.Error.Error()
		return nil, &dbError
	}
	return &gt, nil
}

// // Update takes a group struct and updates specific fields
// func (r *GroupRepository) Update(g models.Group) (*models.Group, error) {
// 	dbc := database.DBCon.Model(&g).Update(map[string]interface{}{
// 		"name":        g.Name,
// 		"description": g.Description,
// 	})
// 	if dbc.Error != nil {
// 		return nil, dbc.Error
// 	}
// 	return &g, nil
// }
//
// // AddGroupMember takes a groupID and a userID and creates a groupMember
// func (r *GroupRepository) AddGroupMember(groupID uint, userID uint) error {
// 	gmr := NewGroupMemberRepository()
// 	var gm models.GroupMember
//
// 	gm.SetDefault()
// 	gm.GroupID = groupID
// 	gm.UserID = userID
//
// 	_, err := gmr.Create(gm)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
