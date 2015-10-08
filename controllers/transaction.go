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
	c.JSON(http.StatusOK, gin.H{"transactions": "Sup"})
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
