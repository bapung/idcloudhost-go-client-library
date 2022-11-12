package vm

import (
	"fmt"
	"regexp"
	"unicode"
)

var validOS = map[string][]string{
	"ubuntu": {"16.04", "18.04", "20.04"},
	"debian": {"9.1"},
	"centos": {"7.3.1611", "6.9.1611"},
}

func validateVmName(name string) error {
	matched, _ := regexp.Match(`^[0-9a-zA-Z][-0-9a-zA-Z]{2,}[0-9a-zA-Z]$`, []byte(name))
	if matched {
		return nil
	}
	return fmt.Errorf("VM validatation failed: VM name must comply regex ^[0-9a-zA-Z][-0-9a-zA-Z]{2,}[0-9a-zA-Z]$")
}

func validateUsername(username string) error {
	matched, _ := regexp.Match(`^[a-zA-Z_][0-9a-zA-Z_-]{1,30}$`, []byte(username))
	if matched {
		return nil
	}
	return fmt.Errorf("VM validatation failed: username must comply regex ^[a-zA-Z_][0-9a-zA-Z_-]{1,30}$")
}

func validatePassword(pass string) error {
	len := 0
	letter := false
	upper := false
	number := false

	for _, c := range pass {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			letter = true
		default:
			//return false, false, false, false
		}
		len++
	}
	if number && upper && letter && len > 7 {
		return nil
	}
	return fmt.Errorf("VM validatation failed: password must contain at least one lowercase and one uppercase ASCII letter (a-z, A-Z) and at least one digit (0-9) and has minimum length of 8 characters")
}

func validateOS(osName string, osVersion string) error {
	for _, v := range validOS[osName] {
		if v == osVersion {
			return nil
		}
	}
	return fmt.Errorf("VM validatation failed: OS %s %s not supported", osName, osVersion)
}

func validateDisks(disks int) error {
	if disks < 20 || disks > 240 {
		return fmt.Errorf("VM validatation failed: ram size must be between 1024 and 65536 MB ")
	}
	return nil
}

func validateRAM(ram int) error {
	if ram < 1024 || ram > 65536 {
		return fmt.Errorf("VM validatation failed: ram size must be between 1024 and 65536 MB ")
	}
	return nil
}

func validateVCPU(vcpu int) error {
	if vcpu < 1 || vcpu > 16 {
		return fmt.Errorf("VM validatation failed: vcpu must be between 1 and 16")
	}
	return nil
}

func validateVmCreateFields(v *NewVM) error {
	if err := validateVmName(v.Name); err != nil {
		return err
	}
	if err := validateVCPU(v.VCPU); err != nil {
		return err
	}
	if err := validateRAM(v.Memory); err != nil {
		return err
	}
	if err := validateDisks(v.Disks); err != nil {
		return err
	}
	if err := validateUsername(v.Username); err != nil {
		return err
	}
	if err := validatePassword(v.InitialPassword); err != nil {
		return err
	}
	if err := validateOS(v.OSName, v.OSVersion); err != nil {
		return err
	}

	return nil
}

func validateVmModifyFields(v *VM) error {
	if v.UUID == "" {
		return fmt.Errorf("UUID is required")
	}
	if v.Name != "" {
		if err := validateVmName(v.Name); err != nil {
			return err
		}
	}
	if err := validateRAM(v.Memory); err != nil {
		return err
	}
	if err := validateVCPU(v.VCPU); err != nil {
		return err
	}
	return nil
}
