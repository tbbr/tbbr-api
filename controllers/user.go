package controllers

import (
	"net/http"
	"payup/database"
	"payup/models"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
)

// UserIndex is used when the user's index is routed to
// this handler will run. Generally, it will
// come with some query parameters like limit and offset
// @returns an array of users
func UserIndex(c *gin.Context) {
	var users []models.User
	database.DBCon.Limit(c.Param("limit")).Find(&users)

	data, err := jsonapi.MarshalToJSON(users)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't marshal to json"})
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}

// UserShow is used to show one specific user
// @returns a user struct
func UserShow(c *gin.Context) {
	var user models.User
	var groups []models.Group
	database.DBCon.First(&user, c.Param("id"))
	database.DBCon.Model(&user).Related(&groups, "Groups")
	user.Groups = groups

	data, err := jsonapi.MarshalToJSON(jsonapi.MarshalIncludedRelations(user))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't marshal to json"})
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}

// UserCreate is used to create one specific user, it'll come with some form data
// @returns the newly created user struct
func UserCreate(c *gin.Context) {

	user := models.User{
		Name:  c.PostForm("name"),
		Email: c.PostForm("email"),
	}

	database.DBCon.Create(&user)

	c.JSON(http.StatusOK, gin.H{"user": user})

}

// UserUpdate is used to update a specific user, it'll also come with some form data
// @returns a user struct
func UserUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"userUpdate": "someContent"})
}

// UserDelete is used to delete one specific user with a `id`
func UserDelete(c *gin.Context) {
	var user models.User
	database.DBCon.First(&user, c.Param("id"))
	database.DBCon.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
