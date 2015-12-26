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
// @returns an array of friendship JSON objects
func FriendshipIndex(c *gin.Context) {
	friendships := []models.Friendship{}
	var curUser models.User
	database.DBCon.First(&curUser, c.Keys["CurrentUserID"])

	database.DBCon.Model(&curUser).Related(&friendships, "Friendships")

	// Get user and friend and friendshipData
	// TODO: n + 1 query problem here, so we'll figure this out later
	for i := range friendships {
		var fd models.FriendshipData
		database.DBCon.First(&friendships[i].Friend, friendships[i].FriendID)
		database.DBCon.First(&fd, friendships[i].FriendshipDataID)

		if curUser.ID == fd.PositiveUserID {
			friendships[i].Balance = fd.Balance
		} else {
			friendships[i].Balance = -fd.Balance
		}

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

// FriendshipShow takes a given ID from gin.Context
// @returns a specific friendship JSON object
func FriendshipShow(c *gin.Context) {
	friendship := models.Friendship{}

	if database.DBCon.First(&friendship, c.Param("id")).RecordNotFound() {
		c.AbortWithError(http.StatusNotFound, appError.RecordNotFound).
			SetMeta(appError.RecordNotFound)
		return
	}

	var fd models.FriendshipData

	database.DBCon.First(&friendship.User, friendship.UserID)
	database.DBCon.First(&friendship.Friend, friendship.FriendID)
	database.DBCon.First(&fd, friendship.FriendshipDataID)

	if friendship.UserID == fd.PositiveUserID {
		friendship.Balance = fd.Balance
	} else {
		friendship.Balance = -fd.Balance
	}

	data, err := jsonapi.MarshalToJSON(friendship)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}
