package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/tbbr/tbbr-api/app-error"
)

// TransactionIndex outputs a certain number of transactions
// TODO: ensure current user is involved in transactions that are returned
func TransactionIndex(c *gin.Context) {
	relatedObjectID := c.Query("relatedObjectId")
	relatedObjectType := c.Query("relatedObjectType")
	isSettledQuery := c.Query("isSettled")
	statusQuery := c.Query("status")
	curUserID := c.Keys["CurrentUserID"]

	var transactions []models.Transaction

	query := database.DBCon

	isSettled, err := strconv.ParseBool(isSettledQuery)
	if isSettledQuery != "" && err == nil {
		query = query.Where("is_settled = ?", isSettled)
	}

	// TODO: Check that statusQuery is a valid status
	if statusQuery != "" {
		query = query.Where("status = ?", statusQuery)
	}

	if relatedObjectID != "" && relatedObjectType != "" {
		query.
			Where("related_object_id = ? AND related_object_type = ?", relatedObjectID, relatedObjectType).
			Order("created_at desc").
			Find(&transactions)
	} else {
		query.
			Where("creator_id = ?", curUserID).
			Find(&transactions)
	}

	// Get creator and relatedUser
	// TODO: n + 1 query problem here, so we'll figure this out later
	for i := range transactions {
		database.DBCon.First(&transactions[i].Recipient, transactions[i].RecipientID)
		database.DBCon.First(&transactions[i].Sender, transactions[i].SenderID)
		database.DBCon.First(&transactions[i].Creator, transactions[i].CreatorID)
	}

	data, err := jsonapi.Marshal(transactions)

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
//		@requires relatedObjectId
//		@requires relatedObjectType
//		@requires recipientId
//		@requires senderid
//		@optional memo
// @returns the newly created transaction
func TransactionCreate(c *gin.Context) {
	var t models.Transaction
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.Unmarshal(buffer, &t)

	if err2 != nil {
		parseFail := appError.JSONParseFailure
		parseFail.Detail = err2.Error()
		c.AbortWithError(http.StatusMethodNotAllowed, err2).
			SetMeta(parseFail)
		return
	}

	t.CreatorID = c.Keys["CurrentUserID"].(uint)

	// Validate our new transaction
	isValid, errApp := t.Validate()

	if isValid == false {
		c.AbortWithError(errApp.Status, errApp).
			SetMeta(errApp)
		return
	}

	database.DBCon.Create(&t)

	database.DBCon.First(&t.Recipient, t.RecipientID)
	database.DBCon.First(&t.Sender, t.SenderID)
	database.DBCon.First(&t.Creator, t.CreatorID)

	data, err := jsonapi.Marshal(&t)

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
//		@optional recipientId
//		@optional senderId
//		@optional amount
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

	err2 := jsonapi.Unmarshal(buffer, &newT)

	if err2 != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	t.Type = newT.Type
	t.Amount = newT.Amount
	t.Memo = newT.Memo
	t.RecipientID = newT.RecipientID
	t.SenderID = newT.SenderID

	// Validate our new transaction
	isValid, errApp := t.Validate()

	if isValid == false {
		c.AbortWithError(errApp.Status, errApp).
			SetMeta(errApp)
		return
	}

	database.DBCon.Save(&t)

	database.DBCon.First(&t.Recipient, t.RecipientID)
	database.DBCon.First(&t.Sender, t.SenderID)
	database.DBCon.First(&t.Creator, t.CreatorID)

	data, err := jsonapi.Marshal(&t)

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
