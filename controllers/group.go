package controllers

import (
	"io/ioutil"
	"net/http"

	"payup/app-error"
	"payup/database"
	"payup/models"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
)

// GroupIndex When the group's index is routed to
// this handler will run. Generally, it will
// come with some query parameters like limit and offset
// @returns an array of group structs
func GroupIndex(c *gin.Context) {
	groups := []models.Group{}
	var curUser models.User
	database.DBCon.First(&curUser, c.Keys["CurrentUserID"])

	database.DBCon.Model(&curUser).Preload("Users").Related(&groups, "Groups")

	data, err := jsonapi.MarshalToJSON(groups)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}

// GroupShow is used to show one specific group, returns a group struct
// @returns a group struct
func GroupShow(c *gin.Context) {
	var group models.Group
	var users []models.User

	if database.DBCon.First(&group, c.Param("id")).RecordNotFound() {
		c.AbortWithError(http.StatusNotFound, appError.RecordNotFound).
			SetMeta(appError.RecordNotFound)
		return
	}

	database.DBCon.Model(&group).Related(&users, "Users")
	group.Users = users
	data, err := jsonapi.MarshalToJSON(group)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}

// GroupCreate is used to create one specific group, it'll come with some form data
// @returns a group struct
func GroupCreate(c *gin.Context) {

	var group models.Group
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.UnmarshalFromJSON(buffer, &group)

	if err2 != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	database.DBCon.Create(&group)

	c.JSON(http.StatusCreated, gin.H{"group": group})
}

// GroupUpdate is used to update a specific group, it'll also come with some form data'
// @returns a group struct
func GroupUpdate(c *gin.Context) {
	var group models.Group
	buffer, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.UnmarshalFromJSON(buffer, &group)

	if err2 != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	// Little hack-ish
	for _, c := range group.UserIDs {
		group.Users = append(group.Users, models.User{ID: c})
	}

	database.DBCon.Model(&group).Association("Users").Replace(group.Users)
	database.DBCon.First(&group, group.ID)
	database.DBCon.Model(&group).Related(&group.Users, "Users")

	data, err := jsonapi.MarshalToJSON(group)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}

// GroupDelete is used to delete one specific group with a `id`
func GroupDelete(c *gin.Context) {
	var group models.Group
	database.DBCon.First(&group, c.Param("id"))
	database.DBCon.Delete(&group)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
