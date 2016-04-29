package fbMessengerBot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"payup/config"

	"github.com/gin-gonic/gin"
)

type fbBotPayload struct {
	Object string  `json:"object"`
	Entry  []entry `json:"entry"`
}

type entry struct {
	ID        string      `json:"id"`
	Time      int         `json:"time"`
	Messaging []messaging `json:"messaging"`
}

type messaging struct {
	Sender    user    `json:"sender"`
	Recipient user    `json:"recipient"`
	Timestamp int     `json:"timestamp"`
	Message   message `json:"message"`
}

type user struct {
	ID string `json:"id"`
}

type message struct {
	MID  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

type sendMessagePayload struct {
	Recipient user    `json:"recipient"`
	Message   message `json:"message"`
}

// VerifyToken function is a route handler
// that is used to verify the token given to from Facebook
func VerifyToken(c *gin.Context) {
	print(c.Query("hub"))
	if c.Query("hub.verify_token") == config.FBMessengerBotToken {
		c.String(200, c.Query("hub.challenge"))
	}
}

// ReceiveMessage function allows the bot to receive a message from a user
// the bot can then subsequently respond to such a message
func ReceiveMessage(c *gin.Context) {
	var payload fbBotPayload

	if c.BindJSON(&payload) == nil {
		messagingEvents := payload.Entry[0].Messaging

		for _, event := range messagingEvents {
			if event.Message.Text == "notify" {
				c.String(200, "OK")
				print("We got here")

				var respMsg sendMessagePayload
				respMsg.Recipient.ID = event.Sender.ID
				respMsg.Message.Text = "Okay, I will start sending you notifications"

				jsonPayload, err := json.Marshal(&respMsg)

				if err != nil {
					http.Post(config.FBMessengerBotPostURL, "application/json", bytes.NewReader(jsonPayload))
				}
			}
		}
	}
	c.String(428, "Couldn't parse json")
}
