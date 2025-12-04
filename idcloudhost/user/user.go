package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	SSHKeys        []SSHKey
	SSHKey         *SSHKey
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

type SSHKey struct {
	ID        int    `json:"id,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type ProfileInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (u *UserAPI) Init(c HTTPClient, authToken string, location string) error {
	u.c = c
	u.AuthToken = authToken
	u.ApiEndpoint = "https://api.idcloudhost.com/v1/user-resource/user"
	u.UserMap = make(map[string]interface{})
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
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}
	return json.NewDecoder(r.Body).Decode(&u.User)
}

func (u *UserAPI) Create() error {
	return NotImplementedError()
}

func (u *UserAPI) ModifyProfile(profileInfo ProfileInfo) error {
	data := url.Values{}
	data.Set("first_name", profileInfo.FirstName)
	data.Set("last_name", profileInfo.LastName)

	req, err := http.NewRequest("PATCH", u.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", u.AuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r, err := u.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&u.User)
}

func (u *UserAPI) ListSSHKeys() error {
	url := fmt.Sprintf("%s/keys", u.ApiEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", u.AuthToken)
	r, err := u.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}
	return json.NewDecoder(r.Body).Decode(&u.SSHKeys)
}

func (u *UserAPI) CreateSSHKey(name string, publicKey string) error {
	url := fmt.Sprintf("%s/keys", u.ApiEndpoint)

	payload := &SSHKey{
		Name:      name,
		PublicKey: publicKey,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", u.AuthToken)

	r, err := u.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&u.SSHKey)
}

func (u *UserAPI) UpdateSSHKeyName(keyID int, name string) error {
	url := fmt.Sprintf("%s/keys/%d", u.ApiEndpoint, keyID)

	payload := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", u.AuthToken)

	r, err := u.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return json.NewDecoder(r.Body).Decode(&u.SSHKey)
}

func (u *UserAPI) DeleteSSHKey(keyID int) error {
	url := fmt.Sprintf("%s/keys/%d", u.ApiEndpoint, keyID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", u.AuthToken)

	r, err := u.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", r.StatusCode))
	}

	return nil
}

func NotImplementedError() error {
	return errors.New("Method not Implemented")
}
