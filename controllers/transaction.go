package controllers

import (
	"io/ioutil"
	"net/http"

	"payup/database"
	"payup/models"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
)

// TransactionIndex outputs a certain number of transactions
// will always be scoped to the current user
func TransactionIndex(c *gin.Context) {
	userID := c.Query("userId")
	groupID := c.Query("groupId")
	curUserID := 11 // Should get this from the Authorization header that gets sent

	var transactions []models.Transaction

	database.DBCon.
		Where("lender_id = ? AND burrower_id = ?", userID, curUserID).
		Or("lender_id = ? AND burrower_id = ?", curUserID, userID).
		Where("group_id = ?", groupID).
		Find(&transactions)

	data, err := jsonapi.MarshalToJSON(transactions)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't marshal to json"})
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)

}

// TransactionCreate will create a transaction that occurs
// between two users in a group
// @parameters
//		@requires	type
//		@requires amount
//		@requires group_id
//		@requires lender_id
//		@requires	burrower_id
// @returns the newly created transaction
func TransactionCreate(c *gin.Context) {
	var transaction models.Transaction
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.UnmarshalFromJSON(buffer, &transaction)

	if err2 != nil {
		c.AbortWithError(http.StatusMethodNotAllowed, err2)
	}

	database.DBCon.Create(&transaction)

	data, err := jsonapi.MarshalToJSON(transaction)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't marshal to json"})
	}

	c.Data(http.StatusCreated, "application/vnd.api+json", data)
}
