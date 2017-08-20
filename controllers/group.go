package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"
	"github.com/tbbr/tbbr-api/repositories"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/tbbr/tbbr-api/app-error"
)

// GroupIndex When the group's index is routed to
// this handler will run. Generally, it will
// come with some query parameters like limit and offset
// @returns an array of group structs
func GroupIndex(c *gin.Context) {
	gr := repositories.NewGroupRepository()
	var curUser models.User
	database.DBCon.First(&curUser, c.Keys["CurrentUserID"])

	data, err := jsonapi.Marshal(gr.List(c.Keys["CurrentUserID"].(uint), 30, 0))

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
	gr := repositories.NewGroupRepository()
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err).
			SetMeta(appError.InvalidParams)
		return
	}
	group, appErr := gr.Get(uint(groupID))
	if appErr != nil {
		c.AbortWithError(http.StatusNotFound, *appErr).SetMeta(*appErr)
		return
	}

	data, err := jsonapi.Marshal(group)
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
	gr := repositories.NewGroupRepository()
	var group models.Group
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.Unmarshal(buffer, &group)

	if err2 != nil {
		parseFail := appError.JSONParseFailure
		parseFail.Detail = err2.Error()
		c.AbortWithError(http.StatusMethodNotAllowed, err2).
			SetMeta(parseFail)
		return
	}

	createdGroup, appErr := gr.Create(group)
	if appErr != nil {
		c.AbortWithError(http.StatusInternalServerError, *appErr).SetMeta(*appErr)
		return
	}
	gr.AddGroupMember(createdGroup.ID, c.Keys["CurrentUserID"].(uint))

	groupWithMembers, appErr2 := gr.Get(createdGroup.ID)
	if appErr != nil {
		c.AbortWithError(http.StatusNotFound, *appErr2).SetMeta(*appErr2)
		return
	}

	data, err3 := jsonapi.Marshal(groupWithMembers)
	if err3 != nil {
		c.AbortWithError(http.StatusInternalServerError, err3).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusCreated, "application/vnd.api+json", data)
}

// GroupUpdate is used to update a specific group
// only updates native fields of the group
// @returns a group struct
func GroupUpdate(c *gin.Context) {
	gr := repositories.NewGroupRepository()
	var group models.Group
	buffer, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	err2 := jsonapi.Unmarshal(buffer, &group)

	if err2 != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	_, appErr := gr.Update(group)
	if appErr != nil {
		c.AbortWithError(http.StatusInternalServerError, *appErr).SetMeta(*appErr)
		return
	}

	groupWithMembers, appErr2 := gr.Get(group.ID)
	if appErr2 != nil {
		c.AbortWithError(http.StatusNotFound, *appErr2).SetMeta(*appErr2)
		return
	}
	data, err := jsonapi.Marshal(groupWithMembers)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}

// GroupJoin handles the request of allowing the current user to join the
// specified group
func GroupJoin(c *gin.Context) {
	gr := repositories.NewGroupRepository()
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err).
			SetMeta(appError.InvalidParams)
		return
	}
	err = gr.AddGroupMember(uint(groupID), c.Keys["CurrentUserID"].(uint))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GroupDelete is used to delete one specific group with a `id`
func GroupDelete(c *gin.Context) {
	var group models.Group
	database.DBCon.First(&group, c.Param("id"))
	database.DBCon.Delete(&group)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
