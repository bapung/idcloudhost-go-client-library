package idcloudhost

import (
	"fmt"
	"regexp"
	"unicode"
)

var validOS = map[string][]string{
	"ubuntu": {"16.04"},
	"debian": {"9.1"},
	"centos": {"7.3.1611", "6.9.1611"},
}

func validateVmName(name string) bool {
	matched, _ := regexp.Match(`^[0-9a-zA-Z][-0-9a-zA-Z]{2,}[0-9a-zA-Z]$`, []byte(name))
	return matched
}

func validateUsername(username string) bool {
	matched, _ := regexp.Match(`^[a-zA-Z_][0-9a-zA-Z_-]{1,30}$`, []byte(username))
	return matched
}

func verifyPassword(pass string) bool {
	len := 0
	letter := true
	upper := false
	special := true
	number := false

	for _, c := range pass {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c):
			letter = true
		default:
			//return false, false, false, false
		}
		len++
	}
	if number && upper && special && letter && len > 7 {
		return true
	}
	return false
}

func validateOS(osName string, osVersion string) bool {
	for _, v := range validOS[osName] {
		if v == osVersion {
			return true
		}
	}
	return false
}

func validateVirtualMachineParam(v map[string]interface{}) error {
	name := v["name"].(string)
	if !validateVmName(name) {
		return fmt.Errorf("VM name must comply regex: ^[0-9a-zA-Z][-0-9a-zA-Z]{2,}[0-9a-zA-Z]$")
	}
	ram := v["ram"].(int)
	if ram < 1024 || ram > 65536 {
		return fmt.Errorf("VM validatation failed: ram size must be between 1024 and 65536 MB ")
	}

	disks := v["disks"].(int)
	if disks < 20 || disks > 240 {
		return fmt.Errorf("VM validatation failed: ram size must be between 1024 and 65536 MB ")
	}

	username := v["username"].(string)
	if !validateUsername(username) {
		return fmt.Errorf("username must comply regex: ^[a-zA-Z_][0-9a-zA-Z_-]{1,30}$")
	}

	password := v["password"].(string)
	if !verifyPassword(password) {
		return fmt.Errorf("password must contain at least one lowercase and one uppercase ASCII letter (a-z, A-Z) and at least one digit (0-9) and has minimum length of 8 characters")
	}

	osName := v["os_name"].(string)
	osVersion := v["os_version"].(string)
	if !validateOS(osName, osVersion) {
		validOSstr := ""
		for k, v := range validOS {
			for _, i := range v {
				validOSstr += fmt.Sprintf("%s %s\n", k, i)
			}
		}
		return fmt.Errorf("OS not supported, currently supported OS are:\n %s", validOSstr)
	}

	return nil
}
