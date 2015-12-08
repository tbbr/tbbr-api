package controllers

import (
	"net/http"

	"payup/app-error"
	"payup/database"
	"payup/models"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
)

// FriendshipIndex takes in query params through
// gin.Context and is restricted to the currentUser
// @returns an array of group structs
func FriendshipIndex(c *gin.Context) {
	friendships := []models.Friendship{}
	var curUser models.User
	database.DBCon.First(&curUser, c.Keys["CurrentUserID"])

	database.DBCon.Model(&curUser).Related(&friendships, "Friendships")

	// Get user and friend
	// TODO: n + 1 query problem here, so we'll figure this out later
	for i := range friendships {
		database.DBCon.First(&friendships[i].Friend, friendships[i].FriendID)
		friendships[i].User = curUser
	}

	data, err := jsonapi.MarshalToJSON(friendships)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}
