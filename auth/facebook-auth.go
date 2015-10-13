package auth

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
	UserID    string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	AvatarURL string
}

type facebookPicture struct {
	PicData data `json:"data"`
}

type data struct {
	URL string `json:"url"`
}

// GetFacebookUserInfo validates an authCode that
// is sent from a client
func GetFacebookUserInfo(authCode string, referrer string) (FacebookUserInfo, error) {
	v := url.Values{}
	v.Set("client_id", os.Getenv("FACEBOOK_KEY"))
	v.Set("client_secret", os.Getenv("FACEBOOK_SECRET"))
	v.Set("redirect_uri", referrer)
	v.Set("code", authCode)

	accessTokenURL := "https://graph.facebook.com/oauth/access_token?" + v.Encode()

	resp, _ := http.Get(accessTokenURL)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	m, _ := url.ParseQuery(string(contents))
	fbAccessToken := m["access_token"][0]

	if fbAccessToken != "" && resp.StatusCode == 200 {
		s := url.Values{}
		s.Set("fields", "id,name,email,gender")
		s.Set("access_token", fbAccessToken)

		resp2, _ := http.Get("https://graph.facebook.com/me?" + s.Encode())

		defer resp2.Body.Close()
		body, _ := ioutil.ReadAll(resp2.Body)

		var userInfo FacebookUserInfo
		json.Unmarshal(body, &userInfo)

		// Get AvatarURL
		picResp, _ := http.Get("https://graph.facebook.com/" + userInfo.UserID + "/picture?type=large&redirect=false")

		defer picResp.Body.Close()
		picBody, _ := ioutil.ReadAll(picResp.Body)

		var fbPic facebookPicture
		json.Unmarshal(picBody, &fbPic)

		// Set userInfo AvatarURL
		if fbPic.PicData.URL != "" {
			userInfo.AvatarURL = fbPic.PicData.URL
		}

		return userInfo, nil
	}

	return FacebookUserInfo{}, errors.New("Failed to get access token")

}
