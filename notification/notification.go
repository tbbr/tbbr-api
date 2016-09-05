package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/tbbr/tbbr-api/database"
	"github.com/tbbr/tbbr-api/models"
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
func New(userID int64) *Notification {
	var deviceToken models.DeviceToken
	// If user doesn't have a device token, return nil
	if database.DBCon.Where("user_id = ?", userID).First(&deviceToken).RecordNotFound() {
		return nil
	}

	var newNotification Notification
	newNotification.To = deviceToken.Token
	newNotification.Priority = "Medium"
	return &newNotification
}

// CreateTransactionTemplate operates on an already created Notification
// it takes a Transaction and
func (n *Notification) CreateTransactionTemplate(t *Transaction) *Notification {
	if t.Sender.Name == "" {
		database.DBCon.First(&t.Sender, t.SenderID)
	}

	if t.Recipient.Name == "" {
		database.DBCon.First(&t.Recipient, t.RecipientID)
	}

	title := fmt.Sprintf("%s's Tab: %s +%s", t.Recipient.Name, t.Sender.Name, t.GetFormattedAmount())
	body := t.Memo
	n.SetDetails(title, body)
}

func (n *Notification) SetDetails(title string, body string) *Notification {
	n.Notification.title = title
	n.Notification.body = body
	return n
}

func (n *Notification) Send() (http.ResponseWriter, error) {
	// If notification title is missing, don't send notification
	if n.Notification.Title == "" {
		return
	}

	data, err := json.Marshal(&n)
	if err != nil {
		fmt.Println("NOTIFICATIONS - failed to marshal: err", err)
		return
	}

	fmt.Println(string(data))

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("NOTIFICATIONS - Couldn't create req for FCM, err: ", err)
	}
	req.Header.Add("Authorization", "key="+os.Getenv("TBBR_FIREBASE_SERVER_KEY"))
	req.Header.Add("Content-Type", "application/json")
	fcmResp, fcmErr := client.Do(req)
	if fcmErr != nil {
		fmt.Println("NOTIFICATIONS - Firebase response failure err: ", fcmErr)
	}
	fmt.Println("NOTIFICATIONS - FCM Response Status ", fcmResp.Status)

	return fcmResp, fcmErr
}
