package controllers

import (
	"net/http"

	"github.com/tbbr/tbbr-api/app-error"
	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"

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

	// database.DBCon.Model(&curUser).Related(&friendships, "Friendships")

	// The ORDER BY clause is a little complex, but the main idea is that we want to show
	// all the friendships that the user owes money to first, then all the friendships
	// that owe the user money, then finally zero balances
	database.DBCon.Raw(`
		SELECT f.* FROM friendships f
		JOIN friendship_data fd 
		ON fd.id = f.friendship_data_id
		WHERE f.user_id = ?
		ORDER BY (
			CASE WHEN (fd.positive_user_id = f.friend_id AND fd.balance > 0) OR (fd.positive_user_id = f.user_id AND fd.balance < 0) then 1 
			WHEN (fd.positive_user_id = f.user_id AND fd.balance > 0) OR (fd.positive_user_id = f.friend_id AND fd.balance < 0) then 2 
			else 3 end
		) ASC, ABS(fd.balance) DESC;
	`, curUser.ID).Scan(&friendships)

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

	data, err := jsonapi.Marshal(friendships)

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

	data, err := jsonapi.Marshal(friendship)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}
