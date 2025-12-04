package api

import (
	"fmt"
	"testing"

	"github.com/bapung/idcloudhost-go-client-library/idcloudhost/vm"
)

const userAuthToken = "h6jyi7lvaqniRk5JhX3FoCExzmh4pkIh"

func TestRiil(t *testing.T) {
	apiC, err := NewClient(userAuthToken, "jkt01")
	if err != nil {
		t.Fatal(fmt.Sprint(err))
	}
	vmC := apiC.VM
	err = vmC.Create(
		vm.NewVM{
			Backup:          false,
			Name:            "testvm",
			OSName:          "ubuntu",
			OSVersion:       "20.04",
			Disks:           20,
			VCPU:            2,
			Memory:          2048,
			Username:        "example",
			InitialPassword: "Password123",
			BillingAccount:  1200132376,
		},
	)
	if err != nil {
		t.Fatal(fmt.Sprint(err))
	}
}
