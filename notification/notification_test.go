package notification

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NotificationTestSuite struct {
	suite.Suite
}

func TestNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}

func (s *NotificationTestSuite) TestNew_Default() {
	testNotif := New("testingToken")
	assert.Equal(s.T(), "testingToken", testNotif.To)
	assert.Equal(s.T(), "high", testNotif.Priority, "The default priority is high")
}

func (s *NotificationTestSuite) TestNew_EmptyToken() {
	testNotif := New("")
	assert.Nil(s.T(), testNotif, "The notification should be nil since EmptyToken was provided")
}

func (s *NotificationTestSuite) TestSetDetails_Default() {
	testNotif := New("testingToken").SetDetails("some_title", "some body")
	assert.Equal(s.T(), "some_title", testNotif.Notification.Title, "Notification title must be set")
	assert.Equal(s.T(), "some body", testNotif.Notification.Body, "Notification body must be set")
}

func (s *NotificationTestSuite) TestSend_EmptyTitle() {
	_, err := New("testingToken").SetDetails("", "some body here").Send()
	assert.NotNil(s.T(), err, "Err must not be nil")
	assert.Equal(s.T(), "Notification title is empty", err.Error())
}

func (s *NotificationTestSuite) TestSend_Default() {
	// Create a new server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var testNotif Notification
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		assert.Nil(s.T(), err, "Error must be nil")
		err = json.Unmarshal(body, &testNotif)
		assert.Nil(s.T(), err, "Error must be nil")

		assert.Equal(s.T(), "testingToken", testNotif.To, "Token is properly set")
		assert.Equal(s.T(), "high", testNotif.Priority, "Priority is high")
		assert.Equal(s.T(), "testing", testNotif.Notification.Title, "Title is set to testing")
		assert.Equal(s.T(), "testing body", testNotif.Notification.Body, "Body is set to testing body")
	}))

	defer os.Setenv("TBBR_FIREBASE_SERVER_URL", os.Getenv("TBBR_FIREBASE_SERVER_URL"))
	defer server.Close()
	os.Setenv("TBBR_FIREBASE_SERVER_URL", server.URL)

	resp, err := New("testingToken").SetDetails("testing", "testing body").Send()

	assert.Nil(s.T(), err, "Error should be nil")
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode, "Response status must be 200 (OK)")
}
