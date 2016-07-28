package models

import (
	"payup/app-error"
	"strconv"
	"time"
)

// DeviceToken model
type DeviceToken struct {
	ID         uint       `json:"id"`
	Token      string     `json:"token"`
	UserID     uint       `json:"userId"`
	DeviceType string     `json:"deviceType"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"-"`
}

type NotificationPayload struct {
	To           string       `json:"to"`
	Priority     string       `json:"priority"`
	Notification Notification `json:"notification"`
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	// Icon  string `json:"icon,omitempty"`
	// Color string `json:"color,omitempty"`
}

// TableName gives gorm information on the name of the table
func (dt DeviceToken) TableName() string {
	return "device_tokens"
}

// Validate the DeviceToken and return a boolean and appError
func (dt DeviceToken) Validate() (bool, appError.Err) {
	if dt.DeviceType != "Android" && dt.DeviceType != "iOS" {
		invalidDeviceType := appError.InvalidParams
		invalidDeviceType.Detail = "The deviceToken deviceType is invalid"
		return false, invalidDeviceType
	}

	if dt.Token == "" {
		invalidToken := appError.InvalidParams
		invalidToken.Detail = "The deviceToken field 'token' cannot be empty"
		return false, invalidToken
	}

	if dt.UserID == 0 {
		invalidUserID := appError.InvalidParams
		invalidUserID.Detail = "The deviceToken userID cannot be 0 or empty"
		return false, invalidUserID
	}

	return true, appError.Err{}
}

////////////////////////////////////////////////////
///////////// API Interface Related ////////////////
////////////////////////////////////////////////////

// GetID returns a stringified version of an ID
func (dt DeviceToken) GetID() string {
	return strconv.FormatUint(uint64(dt.ID), 10)
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (dt *DeviceToken) SetID(id string) error {
	deviceTokenID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	dt.ID = uint(deviceTokenID)
	return nil
}
