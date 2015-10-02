package authr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// FacebookUserInfo represents the model of a user that is returned
// from facebook's oauth
type FacebookUserInfo struct {
	Name   string `json:"name"`
	UserID string `json:"id"`
}

// GetFacebookUserInfo validates an authToken that
// is from facebook client
func GetFacebookUserInfo(authCode string, referrer string) (FacebookUserInfo, error) {
	v := url.Values{}
	v.Set("client_id", os.Getenv("FACEBOOK_KEY"))
	v.Set("client_secret", os.Getenv("FACEBOOK_SECRET"))
	v.Set("redirect_uri", referrer)
	v.Set("code", authCode)

	_ = "breakpoint"

	accessTokenURL := "https://graph.facebook.com/oauth/access_token?" + v.Encode()

	resp, _ := http.Get(accessTokenURL)
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	m, _ := url.ParseQuery(string(contents))
	fbAccessToken := m["access_token"][0]

	if fbAccessToken != "" && resp.StatusCode == 200 {
		// return fbAccessToken
		resp2, _ := http.Get("https://graph.facebook.com/me?access_token=" + fbAccessToken)
		defer resp2.Body.Close()
		body, _ := ioutil.ReadAll(resp2.Body)
		var userInfo FacebookUserInfo
		json.Unmarshal(body, &userInfo)
		return userInfo, nil
	}

	return FacebookUserInfo{}, errors.New("Failed to get access token")

}
