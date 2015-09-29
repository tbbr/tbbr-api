package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// AuthLogin handles the login flow from the user.
func AuthLogin(c *gin.Context) {

	authURL, err := gothic.GetAuthURL(c.Writer, c.Request)

	if err != nil {
		fmt.Printf(err.Error())
	}

	// redirect
	c.Redirect(http.StatusFound, authURL)
}

// AuthCallback is the handler for when a user returns from an auth {provider}
func AuthCallback(c *gin.Context) {

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		fmt.Printf(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
