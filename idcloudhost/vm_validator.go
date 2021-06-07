package idcloudhost

import (
	"fmt"
	"regexp"
)

var validOS = map[string][]string{
	"ubuntu": []string{"16.04"},
	"debian": []string{"9.1"},
	"centos": []string{"7.3.1611", "6.9.1611"},
}

func validateVirutalMachineParam(v map[string]interface{}) error {
	ram := v["ram"].(int)
	if ram < 1024 || ram > 65536 {
		return fmt.Errorf("VM validatation failed: ram size must be between 1024 and 65536 MB ")
	}

	username := v["username"].(string)
	if !validateUsername(username) {
		fmt.Errorf("username must comply regex: ^[a-zA-Z_][0-9a-zA-Z_-]{1,30}$")
	}

	osName := v["os_name"].(string)
	osVersion := v["os_version"].(string)
	if !validateOS(osName, osVersion) {
		fmt.Errorf("OS not supported, currently supported OS are %s")
	}

}

func validateUsername(username string) bool {
	matched, _ := regexp.Match(`^[a-zA-Z_][0-9a-zA-Z_-]{1,30}$`, []byte(username))
	return matched
}

func validateOS(osName string, osVersion string) bool {
	for _, v := range validOS[osName] {
		if v == osVersion {
			return true
		}
	}
	return false
}
