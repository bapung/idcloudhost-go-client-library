package api

import (
	"fmt"
	"testing"

	"github.com/bapung/idcloudhost-go-client-library/idcloudhost/vm"
)

const userAuthToken = "h6jyi7lvaqniRk5JhX3FoCExzmh4pkIh"

// Mock VM client
type mockVMClient struct {
	CreateCalled bool
	CreateArg    vm.NewVM
}

// vmClientMock implements the VM interface
type vmClientMock struct {
	mock *mockVMClient
}

func (v *vmClientMock) Create(newVM vm.NewVM) error {
	v.mock.CreateCalled = true
	v.mock.CreateArg = newVM
	// Simulate success
	return nil
}

func TestRiil(t *testing.T) {
	mockVM := &mockVMClient{}
	vmC := &vmClientMock{mock: mockVM}

	// Simulate API client
	apiC := struct {
		VM interface {
			Create(vm.NewVM) error
		}
	}{
		VM: vmC,
	}

	err := apiC.VM.Create(
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

	if !mockVM.CreateCalled {
		t.Fatal("expected Create to be called on VM client")
	}
}
