package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

// AuthLogin handles the login flow from the user.
func AuthLogin(c *gin.Context) {
	provider, err := gomniauth.Provider(c.Param("provider"))

	if err != nil {
		fmt.Printf(err.Error())
	}

	state := gomniauth.NewState("after", "success")
	authURL, err := provider.GetBeginAuthURL(state, nil)

	if err != nil {
		fmt.Printf(err.Error())
	}

	// redirect
	c.Redirect(http.StatusFound, authURL)
}

// AuthCallback is the handler for when a user returns from an auth {provider}
func AuthCallback(c *gin.Context) {
	provider, err := gomniauth.Provider(c.Param("provider"))

	if err != nil {
		fmt.Printf(err.Error())
	}
	omap, err := objx.FromURLQuery(c.Request.URL.RawQuery)

	if err != nil {
		fmt.Printf(err.Error())
	}

	creds, err := provider.CompleteAuth(omap)

	if err != nil {
		fmt.Printf(err.Error())
	}

	user, err := provider.GetUser(creds)

	if err != nil {
		fmt.Printf(err.Error())
	}

	// _ = "breakpoint"

	c.JSON(http.StatusOK, gin.H{
		"firstName": user.Data().Get("first_name").String(),
		"lastName":  user.Data().Get("last_name").String(),
		"name":      user.Name(),
		"avatarUrl": user.AvatarURL(),
		"email":     user.Email(),
	})
}
