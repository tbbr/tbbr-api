package config

import (
	"fmt"
	"os"
)

// HashIDConfig contains properties of hashId
type HashIDConfig struct {
	Salt      string
	MinLength int
}

var (
	// HashID Config singleton
	HashID = HashIDConfig{"mSwyDdV6Ml4BNvmsM9TK", 11}

	// FBMessengerBotToken holds the verify token needed to verify facebook's
	// messenger bot
	FBMessengerBotToken = "zu4klu2QcPRw64ausbf4"

	// FBMessengerBotPostURL used to send messages to messenger users
	FBMessengerBotPostURL = fmt.Sprintf("https://graph.facebook.com/v2.6/me/messages?access_token=%s", os.Getenv("FB_PAGE_ACCESS_TOKEN"))
)
