package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/tbbr/tbbr-api/app-error"
	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"
)

// DeviceTokenCreate will create a deviceToken for a specified user
// @parameters
//		@requires	deviceType
//		@requires token
// @returns the newly created deviceToken
func DeviceTokenCreate(c *gin.Context) {
	var dt models.DeviceToken
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
	}

	err2 := jsonapi.Unmarshal(buffer, &dt)

	if err2 != nil {
		parseFail := appError.JSONParseFailure
		parseFail.Detail = err2.Error()
		c.AbortWithError(http.StatusMethodNotAllowed, err2).
			SetMeta(parseFail)
		return
	}

	dt.UserID = c.Keys["CurrentUserID"].(uint)

	// Validate our new deviceToken
	isValid, errApp := dt.Validate()

	if isValid == false {
		c.AbortWithError(errApp.Status, errApp).
			SetMeta(errApp)
		return
	}

	database.DBCon.Create(&dt)

	data, err := jsonapi.Marshal(&dt)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).
			SetMeta(appError.JSONParseFailure)
		return
	}

	c.Data(http.StatusCreated, "application/vnd.api+json", data)
}
