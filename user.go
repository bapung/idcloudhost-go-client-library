package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const idcloudhostUserApiUrl = "https://api.idcloudhost.com/v1/user-resource/user"

type UserAPI struct {
	AuthToken string
	User      *User
}

type User struct {
	CookieId     string          `json:"cookie_id"`
	Id           int             `json:"id"`
	LastActivity string          `json:"last_activity"`
	Name         string          `json:"name"`
	ProfileData  UserProfileData `json:"profile_data"`
	SignUpSite   string          `json:"signup_site"`
}

type UserProfileData struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	UpdatedAt string `json:"updatedAt"`
	Avatar    string `json:"avatar"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	CreatedAt string `json:"created_at"`
	Email     string `json:"email"`
}

func (u *UserAPI) setAuthToken(authToken string) {
	u.AuthToken = authToken
}
func (u *UserAPI) getUser() error {
	c := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", idcloudhostUserApiUrl, nil)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("apiKey", u.AuthToken)
	r, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusForbidden {
			return AuthenticationError()
		}
		if r.StatusCode == http.StatusUnauthorized {
			return AuthenticationError()
		}
		return UnknownError()
	}
	return json.NewDecoder(r.Body).Decode(&u.User)
}
