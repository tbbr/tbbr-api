package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"payup/app-error"
	"payup/database"
	"payup/models"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
)

// TransactionIndex outputs a certain number of transactions
// will always be scoped to the current user
func TransactionIndex(c *gin.Context) {
	relatedUserID := c.Query("relatedUserId")
	groupID := c.Query("groupId")
	curUserID := c.Keys["CurrentUserID"]

	var transactions []models.Transaction

	if relatedUserID != "" && groupID != "" {
		database.DBCon.
			Where("related_user_id = ? AND creator_id = ? AND group_id = ?", relatedUserID, curUserID, groupID).
			Or("related_user_id = ? AND creator_id = ? AND group_id = ?", curUserID, relatedUserID, groupID).
			Find(&transactions)
	} else {
		database.DBCon.
			Where("creator_id = ?", curUserID).
			Find(&transactions)
	}

	// Get creator and relatedUser
	// TODO: n + 1 query problem here, so we'll figure this out later
	for i := range transactions {
		database.DBCon.First(&transactions[i].Creator, transactions[i].CreatorID)
		database.DBCon.First(&transactions[i].RelatedUser, transactions[i].RelatedUserID)
	}

	data, err := jsonapi.MarshalToJSON(transactions)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)

}

// TransactionCreate will create a transaction that occurs
// between two users in a group
// @parameters
//		@requires	type
//		@requires amount
//		@requires group_id
//		@requires related_user_id
//		@optional memo
// @returns the newly created transaction along with the updated Balance
func TransactionCreate(c *gin.Context) {
	var t models.Transaction
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.UnmarshalFromJSON(buffer, &t)

	if err2 != nil {
		parseFail := appError.JSONParseFailure
		parseFail.Detail = err2.Error()
		c.AbortWithError(http.StatusMethodNotAllowed, err2).
			SetMeta(parseFail)
		return
	}

	// // Find or create BalanceUser record
	// var bu models.BalanceUser
	// if database.DBCon.
	// 	Where("user_id = ? AND related_user_id = ?", t.LenderID, t.BurrowerID).
	// 	First(&bu).RecordNotFound() {
	// 	// We didn't find a BalanceUser record, so we need to create relationships
	// 	// and create a Balance record as well
	// 	b := models.Balance{
	// 		Amount:         t.Amount,
	// 		PositiveUserID: t.LenderID,
	// 	}
	//
	// 	database.DBCon.Create(&b)
	//
	// 	bu := models.BalanceUser{
	// 		BalanceID:     b.ID,
	// 		UserID:        t.LenderID,
	// 		RelatedUserID: t.BurrowerID,
	// 		GroupID:       t.GroupID,
	// 	}
	//
	// 	bu2 := models.BalanceUser{
	// 		BalanceID:     b.ID,
	// 		UserID:        t.BurrowerID,
	// 		RelatedUserID: t.LenderID,
	// 		GroupID:       t.GroupID,
	// 	}
	//
	// 	database.DBCon.Create(&bu)
	// 	database.DBCon.Create(&bu2)
	//
	// } else {
	// 	// Update Balance with the newly created transaction
	// }

	// // Attach BalanceID to Transaction
	// t.BalanceID = bu.BalanceID

	t.CreatorID = c.Keys["CurrentUserID"].(uint)

	fmt.Print(t)

	database.DBCon.Create(&t)

	database.DBCon.First(&t.RelatedUser, t.RelatedUserID)
	database.DBCon.First(&t.Creator, t.CreatorID)

	data, err := jsonapi.MarshalToJSON(&t)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusCreated, "application/vnd.api+json", data)
}
