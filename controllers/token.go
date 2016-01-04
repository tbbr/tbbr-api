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
	grantType := c.PostForm("grant_type")
	accessToken := c.PostForm("access_token")
	var err error

	if grantType == "facebook_auth_code" || accessToken == "" {
		accessToken, err = auth.GetFacebookAccessToken(c.PostForm("auth_code"), c.Request.Referer())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	userInfo, err2 := auth.GetFacebookUserInfo(accessToken)

	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
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

	auth.UpdateFacebookUserFriends(accessToken, user)

	token := models.Token{
		Category: "oAuth",
		UserID:   user.ID,
	}

	database.DBCon.Create(&token)

	c.JSON(http.StatusOK, token)
}
