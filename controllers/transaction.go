package controllers

import (
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
// @returns the newly created transaction
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

	t.CreatorID = c.Keys["CurrentUserID"].(uint)

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

// TransactionUpdate will update an existing transaction
// between two users in a group
// @parameters
//		@requires id
//		@optional	type
//		@optional amount
//		@optional group_id
//		@optional related_user_id
//		@optional memo
// @returns the updated transaction
func TransactionUpdate(c *gin.Context) {
	var t models.Transaction
	var newT models.Transaction

	if database.DBCon.First(&t, c.Param("id")).RecordNotFound() {
		c.AbortWithError(http.StatusNotFound, appError.RecordNotFound).
			SetMeta(appError.RecordNotFound)
		return
	}

	// Ensure current user is creator of transaction
	if t.CreatorID != c.Keys["CurrentUserID"].(uint) {
		c.AbortWithError(appError.InsufficientPermission.Status, appError.InsufficientPermission).
			SetMeta(appError.InsufficientPermission)
		return
	}

	buffer, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.UnmarshalFromJSON(buffer, &newT)

	if err2 != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	t.Type = newT.Type
	t.Amount = newT.Amount
	t.Memo = newT.Memo

	database.DBCon.Save(&t)

	database.DBCon.First(&t.RelatedUser, t.RelatedUserID)
	database.DBCon.First(&t.Creator, t.CreatorID)

	data, err := jsonapi.MarshalToJSON(&t)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)

}

// TransactionDelete will delete an existing transaction
// or throw an error
// @parameters
//		@requires id
// @returns JSON meta property with status
func TransactionDelete(c *gin.Context) {
	var t models.Transaction
	if database.DBCon.First(&t, c.Param("id")).RecordNotFound() {
		c.AbortWithError(http.StatusNotFound, appError.RecordNotFound).
			SetMeta(appError.RecordNotFound)
		return
	}

	// Ensure current user is creator of transaction
	if t.CreatorID != c.Keys["CurrentUserID"].(uint) {
		c.AbortWithError(appError.InsufficientPermission.Status, appError.InsufficientPermission).
			SetMeta(appError.InsufficientPermission)
		return
	}

	database.DBCon.Delete(&t)

	c.JSON(http.StatusOK, gin.H{"meta": gin.H{"success": true}})
}
