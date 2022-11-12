package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type UserAPI struct {
	c              HTTPClient
	AuthToken      string
	ApiEndpoint    string
	BillingAccount []string
	User           *User
	UserMap        map[string]interface{}
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

func (u *UserAPI) Init(c HTTPClient, authToken string, location string) error {
	u.c = c
	u.AuthToken = authToken
	u.ApiEndpoint = "https://api.idcloudhost.com/v1/user-resource/user"
	return nil
}

func (u *UserAPI) Get(uuid string) error {
	req, err := http.NewRequest("GET", u.ApiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", u.AuthToken)
	r, err := u.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&u.User)
}

func (u *UserAPI) Create() error {
	return NotImplementedError()
}

func (u *UserAPI) Modify() error {
	return NotImplementedError()
}

func (u *UserAPI) Delete() error {
	return NotImplementedError()
}
