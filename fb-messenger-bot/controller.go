package fbMessengerBot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"payup/config"

	"github.com/gin-gonic/gin"
)

type fbBotPayload struct {
	Object string  `json:"object"`
	Entry  []entry `json:"entry"`
}

type entry struct {
	ID        int         `json:"id"`
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
	ID int `json:"id"`
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
	// foundNotifyMsg := false

	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		print(err.Error())
		c.String(500, "Couldn't parse")
	}

	marshalErr := json.Unmarshal(body, &payload)

	if marshalErr == nil {

		print(payload.Object)
		print("\n")
		print(len(payload.Entry))
		print("\n")

		messagingEvents := payload.Entry[0].Messaging
		for _, event := range messagingEvents {
			print(event.Sender.ID)
			if event.Message.Text == "notify" {
				print("user wants to be notified")
				var respMsg sendMessagePayload
				respMsg.Recipient.ID = event.Sender.ID
				respMsg.Message.Text = "Okay, I will start sending you notifications"

				jsonPayload, err := json.Marshal(&respMsg)

				if err == nil {
					resp, err := http.Post(config.FBMessengerBotPostURL, "application/json", bytes.NewReader(jsonPayload))

					if err == nil {
						print(resp.StatusCode)
					} else {
						print(err.Error())
					}

				} else {
					print(err.Error())
				}
			} else {
				print(event.Message.Text)
			}
		}
		c.String(200, payload.Object)
	} else {
		print(marshalErr.Error())
		c.String(428, "Couldn't parse json")
		return
	}

	// if c.BindJSON(&payload) == nil {
	// 	messagingEvents := payload.Entry[0].Messaging
	//
	// 	for _, event := range messagingEvents {
	// 		if event.Message.Text == "notify" {
	// 			foundNotifyMsg = true
	// 			print("We got here")
	//
	// 			var respMsg sendMessagePayload
	// 			respMsg.Recipient.ID = event.Sender.ID
	// 			respMsg.Message.Text = "Okay, I will start sending you notifications"
	//
	// 			jsonPayload, err := json.Marshal(&respMsg)
	//
	// 			if err != nil {
	// 				http.Post(config.FBMessengerBotPostURL, "application/json", bytes.NewReader(jsonPayload))
	// 			}
	// 		} else {
	// 			var respMsg sendMessagePayload
	// 			respMsg.Recipient.ID = event.Sender.ID
	// 			respMsg.Message.Text = "You messaged me!"
	//
	// 			jsonPayload, err := json.Marshal(&respMsg)
	//
	// 			if err != nil {
	// 				http.Post(config.FBMessengerBotPostURL, "application/json", bytes.NewReader(jsonPayload))
	// 			}
	// 		}
	// 	}
	// } else {
	// 	c.String(428, "Couldn't parse json")
	// 	return
	// }

	// if foundNotifyMsg == true {
	// 	c.String(200, "Notified")
	// } else {
	// 	c.String(200, "Worked")
	// }
}
