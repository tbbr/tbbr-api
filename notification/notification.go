package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Notification struct {
	To           string              `json:"to"`
	Priority     string              `json:"priority"`
	Notification NotificationDetails `json:"notification"`
	// Icon  string `json:"icon,omitempty"`
	// Color string `json:"color,omitempty"`
}

type NotificationDetails struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// New takes in a userID and returns a Notification
func New(token string) *Notification {
	var newNotification Notification
	newNotification.To = token
	newNotification.Priority = "Medium"
	return &newNotification
}

// SetDetails sets the title and body of the notification
func (n *Notification) SetDetails(title string, body string) *Notification {
	n.Notification.Title = title
	n.Notification.Body = body
	return n
}

// Send method will send the notification through FCM and return the response
func (n *Notification) Send() (*http.Response, error) {
	// If notification title is missing, don't send notification
	if n.Notification.Title == "" {
		return nil, errors.New("Notification title is empty")
	}
	data, err := json.Marshal(&n)
	if err != nil {
		fmt.Println("NOTIFICATIONS - failed to marshal: err", err)
		return nil, err
	}
	fmt.Println(string(data))
	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("NOTIFICATIONS - Couldn't create req for FCM, err: ", err)
	}
	req.Header.Add("Authorization", "key="+os.Getenv("TBBR_FIREBASE_SERVER_KEY"))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	fcmResp, fcmErr := client.Do(req)
	if fcmErr != nil {
		fmt.Println("NOTIFICATIONS - Firebase response failure err: ", fcmErr)
	}
	fmt.Println("NOTIFICATIONS - FCM Response Status ", fcmResp.Status)

	return fcmResp, fcmErr
}
