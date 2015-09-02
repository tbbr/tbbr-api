package controllers

import (
	"net/http"

	"payup/models"
	"payup/database"
	"github.com/gin-gonic/gin"
)

// TransactionIndex outputs a certain number of transactions
// will always be scoped to the current user
func TransactionIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"transactions": "Sup"})
}

// TransactionCreate will create a transaction that occurs
// between two users in a group
// @returns the newly created transaction
func TransactionCreate(c *gin.Context) {
	var t models.Transaction
	t.Amount := 150
	t.Comment := "I bought Brandon Coffee"
	t.LenderID := 1
	t.BurrowerID := 3
	t.GroupID := 2

	database.DBCon.Create(&t)

	c.JSON(http.StatusOK, gin.H{"transaction": t})
}
