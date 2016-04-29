package fbMessengerBot

import (
	"payup/config"

	"github.com/gin-gonic/gin"
)

// VerifyToken function is a route handler
// that is used to verify the token given to from Facebook
func VerifyToken(c *gin.Context) {
	if c.Query("verify_token") == config.FBMessengerBotToken {
		c.JSON(200, gin.H{
			"challenge": c.Query("challenge"),
		})
	}
}

// ReceiveMessage function allows the bot to receive a message from a user
// the bot can then subsequently respond to such a message
func ReceiveMessage(c *gin.Context) {

}