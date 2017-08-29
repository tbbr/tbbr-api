package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"
	"github.com/tbbr/tbbr-api/repositories"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/tbbr/tbbr-api/app-error"
)

// GroupTransactionIndex outputs a certain number of groupTransactions
// will always be scoped to the current user
func GroupTransactionIndex(c *gin.Context) {
	gtr := repositories.NewGroupTransactionRepository()
	groupID, err := strconv.ParseUint(c.Query("groupId"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err).
			SetMeta(appError.InvalidParams)
		return
	}
	data, err := jsonapi.Marshal(gtr.List(uint(groupID), 30, 0))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}

// GroupTransactionCreate will create a groupTransaction
// @parameters
//		@requires amount
//		@requires senders
//		@requires recipients
//		@requires senderSplits
//		@requires recipientSplits
//		@requires senderSplitType
//    @requires recipientSplitType
//    @requires groupID
//		@optional memo
// @returns the newly created transaction
func GroupTransactionCreate(c *gin.Context) {
	gtr := repositories.NewGroupTransactionRepository()
	var gt models.GroupTransaction
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.Unmarshal(buffer, &gt)
	if err2 != nil {
		parseFail := appError.JSONParseFailure
		parseFail.Detail = err2.Error()
		c.AbortWithError(http.StatusMethodNotAllowed, err2).
			SetMeta(parseFail)
		return
	}

	gt.CreatorID = c.Keys["CurrentUserID"].(uint)

	createdGt, appErr := gtr.Create(gt)
	if appErr != nil {
		c.AbortWithError(http.StatusInternalServerError, *appErr).SetMeta(*appErr)
		return
	}

	gtNew, appErr2 := gtr.Get(createdGt.ID)
	if appErr2 != nil {
		c.AbortWithError(http.StatusNotFound, *appErr2).SetMeta(*appErr2)
		return
	}

	data, err := jsonapi.Marshal(gtNew)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusCreated, "application/vnd.api+json", data)
}

// // TransactionUpdate will update an existing transaction
// // between two users in a group
// // @parameters
// //		@requires id
// //		@optional	type
// //		@optional recipientId
// //		@optional senderId
// //		@optional amount
// //		@optional memo
// // @returns the updated transaction
// func TransactionUpdate(c *gin.Context) {
// 	var t models.Transaction
// 	var newT models.Transaction
//
// 	if database.DBCon.First(&t, c.Param("id")).RecordNotFound() {
// 		c.AbortWithError(http.StatusNotFound, appError.RecordNotFound).
// 			SetMeta(appError.RecordNotFound)
// 		return
// 	}
//
// 	// Ensure current user is creator of transaction
// 	if t.CreatorID != c.Keys["CurrentUserID"].(uint) {
// 		c.AbortWithError(appError.InsufficientPermission.Status, appError.InsufficientPermission).
// 			SetMeta(appError.InsufficientPermission)
// 		return
// 	}
//
// 	buffer, err := ioutil.ReadAll(c.Request.Body)
//
// 	if err != nil {
// 		c.AbortWithError(http.StatusNotAcceptable, err)
// 	}
//
// 	err2 := jsonapi.Unmarshal(buffer, &newT)
//
// 	if err2 != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err).
// 			SetMeta(appError.JSONParseFailure)
// 		return
// 	}
//
// 	t.Type = newT.Type
// 	t.Amount = newT.Amount
// 	t.Memo = newT.Memo
// 	t.RecipientID = newT.RecipientID
// 	t.SenderID = newT.SenderID
//
// 	// Validate our new transaction
// 	isValid, errApp := t.Validate()
//
// 	if isValid == false {
// 		c.AbortWithError(errApp.Status, errApp).
// 			SetMeta(errApp)
// 		return
// 	}
//
// 	database.DBCon.Save(&t)
//
// 	database.DBCon.First(&t.Recipient, t.RecipientID)
// 	database.DBCon.First(&t.Sender, t.SenderID)
// 	database.DBCon.First(&t.Creator, t.CreatorID)
//
// 	data, err := jsonapi.Marshal(&t)
//
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err).
// 			SetMeta(appError.JSONParseFailure)
// 		return
// 	}
//
// 	c.Data(http.StatusOK, "application/vnd.api+json", data)
// }
//
// GroupTransactionDelete will delete an existing groupTransaction
// or throw an error
// @parameters
//		@requires id
// @returns JSON meta property with status
func GroupTransactionDelete(c *gin.Context) {
	var gt models.GroupTransaction
	if database.DBCon.First(&gt, c.Param("id")).RecordNotFound() {
		c.AbortWithError(http.StatusNotFound, appError.RecordNotFound).
			SetMeta(appError.RecordNotFound)
		return
	}

	// Ensure current user is creator of transaction
	if gt.CreatorID != c.Keys["CurrentUserID"].(uint) {
		c.AbortWithError(appError.InsufficientPermission.Status, appError.InsufficientPermission).
			SetMeta(appError.InsufficientPermission)
		return
	}

	database.DBCon.Delete(&gt)

	c.JSON(http.StatusOK, gin.H{"meta": gin.H{"success": true}})
}
