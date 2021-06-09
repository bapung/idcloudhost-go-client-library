package idcloudhost

import (
	"fmt"
	"log"
	"testing"
)

func TestNonValidName(t *testing.T) {
	NonValidName := []string{
		"ax1", "-startwhyphen", "endwhypen-", "contains@symbol", "contains_underscore",
	}
	for _, name := range NonValidName {
		if validateVmName(name) {
			t.Fatal(fmt.Errorf("validate VM name %s should return False", name))
		}
	}
}

func TestValidOS(t *testing.T) {
	ValidOSes := map[string][]string{
		"ubuntu": {"16.04"},
		"debian": {"9.1"},
		"centos": {"7.3.1611", "6.9.1611"},
	}

	for k, v := range ValidOSes {
		for _, i := range v {
			if !validateOS(k, i) {
				t.Fatal(fmt.Errorf("validate OS %s %s should return true", k, v))
			}
		}
	}
}

func TestNonValidOS(t *testing.T) {
	ValidOSes := map[string][]string{
		"manjaro": {"iamrolling"},
	}

	for k, v := range ValidOSes {
		for _, i := range v {
			if validateOS(k, i) {
				t.Fatal(fmt.Errorf("validate OS %s %s should return false", k, v))
			}
		}
	}
}

func TestNonValidPassword(t *testing.T) {
	PasswordAndErr := map[string]string{
		"aA123":        "password length < 8 characters",
		"abcdefg12345": "password does not contain at least 1 uppercase character",
		"abcdEFGHIJK":  "password does not contain at least 1 digit",
		"ABCDEF12345":  "password does not contain at least 1 lowercase character",
	}
	for pass, errStr := range PasswordAndErr {
		if validatePassword(pass) {
			t.Fatal(fmt.Errorf("this password is valid: %s; should not valid in this test", pass))
		}
		log.Printf("expected error: %s for password: %s", errStr, pass)
	}
}
