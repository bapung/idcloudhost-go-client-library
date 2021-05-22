package idcloudhost

import (
	"log"
	"testing"
)

func TestGetUser(t *testing.T) {
	u := UserAPI{}
	u.Init(userAuthToken)
	if err := u.Get(); err != nil {
		t.Fatal(err)
	}
	log.Println(u.User)
}
func TestGetUserNotAuthorized(t *testing.T) {
	u := UserAPI{}
	u.Init("non-valid-auth-token")
	if err := u.Get(); err == nil {
		t.Fatal(err)
	}
}
