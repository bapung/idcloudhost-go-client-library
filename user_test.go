package main

import (
	"log"
	"testing"
)

const userAuthToken = ""

func TestGetUser(t *testing.T) {
	u := UserAPI{}
	u.setAuthToken(userAuthToken)
	if err := u.getUser(); err != nil {
		t.Fatal(err)
	}
	log.Println(u.User)
}
func TestGetUserNotAuthorized(t *testing.T) {
	u := UserAPI{}
	u.setAuthToken("non-valid-auth-token")
	if err := u.getUser(); err == nil {
		t.Fatal(err)
	}
}
