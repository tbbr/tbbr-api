package controllers

import (
	"net/http"
	"payup/auth"
	"payup/database"
	"payup/models"

	"github.com/gin-gonic/gin"
)

// TokenOAuthGrant grants an oAuth token and creates a user if not
// created already
func TokenOAuthGrant(c *gin.Context) {
	userInfo, err := auth.GetFacebookUserInfo(c.PostForm("auth_code"), c.Request.Referer())

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
	}

	// Find or create user
	var user models.User

	database.DBCon.Where(models.User{
		ExternalID: userInfo.UserID,
	}).Attrs(models.User{
		Name:      userInfo.Name,
		Email:     userInfo.Email,
		Gender:    userInfo.Gender,
		AvatarURL: userInfo.AvatarURL,
	}).FirstOrCreate(&user)

	auth.UpdateFacebookUserFriends(userInfo.AccessToken, user)

	token := models.Token{
		Category: "oAuth",
		UserID:   user.ID,
	}

	database.DBCon.Create(&token)

	c.JSON(http.StatusOK, token)
}
