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

// GroupIndex When the groupMember's index is routed to
// this handler will run. Generally, it will
// come with some query parameters like group_id
// @returns an array of groupMember structs
func GroupMemberIndex(c *gin.Context) {
	gm := repositories.NewGroupMemberRepository()
	groupID, err := strconv.ParseUint(c.Query("groupId"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err).
			SetMeta(appError.InvalidParams)
		return
	}

	data, err := jsonapi.Marshal(gm.List(uint(groupID), 30, 0))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusOK, "application/vnd.api+json", data)
}

// GroupMemberCreate is used to create a group member
// @returns a group struct
func GroupMemberCreate(c *gin.Context) {
	gm := repositories.NewGroupMemberRepository()
	var groupMember models.GroupMember
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.Unmarshal(buffer, &groupMember)

	if err2 != nil {
		parseFail := appError.JSONParseFailure
		parseFail.Detail = err2.Error()
		c.AbortWithError(http.StatusMethodNotAllowed, err2).
			SetMeta(parseFail)
		return
	}

	createdGroupMember, appErr := gm.Create(groupMember)
	if appErr != nil {
		c.AbortWithError(http.StatusInternalServerError, *appErr).SetMeta(*appErr)
		return
	}

	groupMemberWithUser, appErr2 := gm.Get(createdGroupMember.ID)
	if appErr != nil {
		c.AbortWithError(http.StatusNotFound, *appErr2).SetMeta(*appErr2)
		return
	}

	data, err3 := jsonapi.Marshal(groupMemberWithUser)
	if err3 != nil {
		c.AbortWithError(http.StatusInternalServerError, err3).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusCreated, "application/vnd.api+json", data)
}

// GroupMemberDelete is used to delete one specific group with a `id`
func GroupMemberDelete(c *gin.Context) {
	var groupMember models.GroupMember
	database.DBCon.First(&groupMember, c.Param("id"))
	database.DBCon.Delete(&groupMember)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
