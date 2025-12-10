//go:build !integration

package vm

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

type HTTPClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (H HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	return H.DoFunc(r)
}

const userAuthToken = "h6jyi7lvaqniRk5JhX3FoCExzmh4pkIh"

var (
	mockHttpClient = &HTTPClientMock{}
	testVmApi      = VirtualMachineAPI{}
	loc            = "jkt01"
)

func TestRiil(t *testing.T) {
	if err := testVmApi.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize vm api: %v", err)
	}
}

func TestGetVMbyUUID(t *testing.T) {
	if err := testVmApi.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize vm api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "this-is-a-supposed-to-be-valid",
			},
			Body:       `{"backup":true,"billing_account":6,"created_at":"2018-02-22 11:10:17","description":"","hostname":"hostname","hypervisor_id":null,"id":7,"mac":"52:54:00:6c:6a:ac","memory":2048,"name":"Ubuntu-16-04","os_name":"ubuntu","os_version":"16.04","private_ipv4":"","status":"running","storage":[{"created_at":"2018-02-22 11:10:37.793878","id":5,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a"}],"tags":null,"updated_at":"2018-02-22 13:48:21","user_id":8,"username":"example","uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a","vcpu":2}`,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := testVmApi.Get(test.RequestData["uuid"].(string))
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		log.Println(testVmApi.VM)
	}
}
func TestModify(t *testing.T) {
	if err := testVmApi.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize vm api: %v", err)
	}
	testCases := []struct {
		RequestData VM
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: VM{
				UUID:   "this-is-a-supposed-to-be-valid",
				Name:   "updated-vm-name",
				VCPU:   1,
				Memory: 1024,
			},
			Body:       `{"backup":true,"billing_account":6,"created_at":"2018-02-22 11:10:17","description":"","hostname":"hostname","hypervisor_id":null,"id":7,"mac":"52:54:00:6c:6a:ac","memory":2048,"name":"Ubuntu-16-04","os_name":"ubuntu","os_version":"16.04","private_ipv4":"","status":"running","storage":[{"created_at":"2018-02-22 11:10:37.793878","id":5,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a"}],"tags":null,"updated_at":"2018-02-22 13:48:21","user_id":8,"username":"example","uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a","vcpu":2}`,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
		{
			RequestData: VM{
				UUID: "this-is-a-supposed-to-be-valid",
				Name: "__name-notvalid-vm",
			},
			Body:       ``,
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf(`VM validatation failed: VM name must comply regex ^[0-9a-zA-Z][-0-9a-zA-Z]{2,}[0-9a-zA-Z]$`),
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := testVmApi.Modify(test.RequestData)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}

func TestListAllVMs(t *testing.T) {
	if err := testVmApi.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize vm api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: nil,
			Body:        `[{"backup":false,"billing_account":6,"created_at":"2018-02-22 14:24:30","description":"","hostname":"hostname","hypervisor_id":null,"id":11,"mac":"52:54:00:59:44:d1","memory":2048,"name":"Name of the Clone","os_name":"ubuntu","os_version":"16.04","private_ipv4":"10.1.14.251","status":"running","storage":[{"created_at":"2018-02-22 14:24:30.312877","id":9,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"d582f16a-013b-4a23-8463-c66bbbc96c43"}],"tags":null,"updated_at":null,"user_id":8,"username":"example","uuid":"d582f16a-013b-4a23-8463-c66bbbc96c43","vcpu":2},{"backup":false,"billing_account":6,"created_at":"2018-02-22 14:24:03","description":"","hostname":"hostname","hypervisor_id":null,"id":10,"mac":"52:54:00:a2:52:6a","memory":2048,"name":"Ubuntu-16-04","os_name":"ubuntu","os_version":"16.04","private_ipv4":"10.1.14.253","status":"running","storage":[{"created_at":"2018-02-22 14:24:13.766985","id":8,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"fc880f74-cf03-4a7a-93da-74c506157023"}],"tags":null,"updated_at":"2018-02-22 14:24:13","user_id":8,"username":"example","uuid":"fc880f74-cf03-4a7a-93da-74c506157023","vcpu":2}]`,
			StatusCode:  http.StatusOK,
			Error:       nil,
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := testVmApi.ListAll()
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		log.Println(testVmApi.VMList)
	}
}
func TestCreateVM(t *testing.T) {
	if err := testVmApi.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize vm api: %v", err)
	}
	testCases := []struct {
		RequestData NewVM
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: NewVM{
				Backup:          false,
				Name:            "testvm",
				OSName:          "ubuntu",
				OSVersion:       "16.04",
				Disks:           20,
				VCPU:            1,
				Memory:          1024,
				Username:        "example",
				InitialPassword: "Password123",
				BillingAccount:  9999,
			},
			Body:       `{"backup":true,"billing_account":6,"created_at":"2018-02-22 11:10:17","description":"","hostname":"hostname","hypervisor_id":null,"id":7,"mac":"52:54:00:6c:6a:ac","memory":2048,"name":"Ubuntu-16-04","os_name":"ubuntu","os_version":"16.04","private_ipv4":"","status":"running","storage":[{"created_at":"2018-02-22 11:10:37.793878","id":5,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a"}],"tags":null,"updated_at":"2018-02-22 13:48:21","user_id":8,"username":"example","uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a","vcpu":2}`,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
		{
			RequestData: NewVM{
				Name: "incomplete-vm",
			},
			Body:       ``,
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf(`VM validatation failed: vcpu must be between 1 and 16`),
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := testVmApi.Create(test.RequestData)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		log.Print(testVmApi.VM)
	}
}

func TestToggleBackup(t *testing.T) {
	if err := testVmApi.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize vm api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]string
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]string{"uuid": "this-is-valid-uuid"},
			Body:        `{"backup":true,"billing_account":6,"created_at":"2018-02-22 11:10:17","description":"","hostname":"hostname","hypervisor_id":null,"id":7,"mac":"52:54:00:6c:6a:ac","memory":2048,"name":"Ubuntu-16-04","os_name":"ubuntu","os_version":"16.04","private_ipv4":"","status":"running","storage":[{"created_at":"2018-02-22 11:10:37.793878","id":5,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a"}],"tags":null,"updated_at":"2018-02-22 13:48:21","user_id":8,"username":"example","uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a","vcpu":2}`,
			StatusCode:  http.StatusOK,
			Error:       nil,
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := testVmApi.ToggleAutoBackup(test.RequestData["uuid"])
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		log.Print(testVmApi.VM)
	}
}
func TestClone(t *testing.T) {
	if err := testVmApi.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize vm api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]string
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]string{"uuid": "this-is-valid-uuid", "clone_name": "test-clone-name"},
			Body:        `{"backup":true,"billing_account":6,"created_at":"2018-02-22 11:10:17","description":"","hostname":"hostname","hypervisor_id":null,"id":7,"mac":"52:54:00:6c:6a:ac","memory":2048,"name":"Ubuntu-16-04","os_name":"ubuntu","os_version":"16.04","private_ipv4":"","status":"running","storage":[{"created_at":"2018-02-22 11:10:37.793878","id":5,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a"}],"tags":null,"updated_at":"2018-02-22 13:48:21","user_id":8,"username":"example","uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a","vcpu":2}`,
			StatusCode:  http.StatusOK,
			Error:       nil,
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := testVmApi.Clone(test.RequestData["uuid"], test.RequestData["clone_name"])
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		log.Print(testVmApi.VM)
	}
}
