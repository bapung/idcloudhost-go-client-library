package vm

import (
	"fmt"
	"log"
	"testing"
)

func TestInvalidVMName(t *testing.T) {
	invalidName := []string{
		"ax1", "-startwhyphen", "endwhypen-", "contains@symbol", "contains_underscore",
	}
	for _, name := range invalidName {
		if err := validateVmName(name); err == nil {
			t.Fatal(fmt.Errorf("validate VM name %s should return False", name))
		}
	}
}

func TestValidOS(t *testing.T) {
	validOSes := map[string][]string{
		"ubuntu": {"16.04"},
		"debian": {"9.1"},
		"centos": {"7.3.1611", "6.9.1611"},
	}

	for k, v := range validOSes {
		for _, i := range v {
			if err := validateOS(k, i); err != nil {
				t.Fatal(fmt.Errorf("validate OS %s %s should return true", k, v))
			}
		}
	}
}

func TestInvalidOS(t *testing.T) {
	validOSes := map[string][]string{
		"manjaro": {"iamrolling"},
	}

	for k, v := range validOSes {
		for _, i := range v {
			if err := validateOS(k, i); err == nil {
				t.Fatal(fmt.Errorf("validate OS %s %s should return error", k, v))
			}
		}
	}
}

func TestInvalidUsername(t *testing.T) {
	invalidName := []string{
		"+startwithnonchara", "00startwithnumber", "contains@symbol", "thisusernameeeeeeeeeistoootooolooooooonnnnggggg",
	}
	for _, name := range invalidName {
		if err := validateUsername(name); err == nil {
			t.Fatal(fmt.Errorf("validate VM name %s should return error", name))
		}
	}
}

func TestInvalidPassword(t *testing.T) {
	passwordAndErr := map[string]string{
		"aA123":        "password length < 8 characters",
		"abcdefg12345": "password does not contain at least 1 uppercase character",
		"abcdEFGHIJK":  "password does not contain at least 1 digit",
		"ABCDEF12345":  "password does not contain at least 1 lowercase character",
	}
	for pass, errStr := range passwordAndErr {
		if err := validatePassword(pass); err == nil {
			t.Fatal(fmt.Errorf("this password is valid: %s; should not valid in this test", pass))
		}
		log.Printf("expected error: %s for password: %s", errStr, pass)
	}
}

func TestInvalidRAM(t *testing.T) {
	invalidRAM := []int{-1, 0, 200}
	for _, ram := range invalidRAM {
		if err := validateRAM(ram); err == nil {
			t.Fatal(fmt.Errorf("RAM value is valid %d; should not valid for this test", ram))
		}
	}
}

func TestInvalidVCPU(t *testing.T) {
	invalidVCPUs := []int{-2, 0, 200}
	for _, vcpu := range invalidVCPUs {
		if err := validateVCPU(vcpu); err == nil {
			t.Fatal(fmt.Errorf("VCPU value is valid %d; should not valid for this test", vcpu))
		}
	}
}
