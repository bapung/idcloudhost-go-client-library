package idcloudhost

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

var testDiskApi DiskAPI

func TestGetDiskByUUID(t *testing.T) {
	testDiskApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"vm_uuid": "valid-vm-uuid",
				"uuid":    "this-is-NON-exist-disk-uuid",
			},
			Body:       `{"backup":true,"billing_account":6,"created_at":"2018-02-22 11:10:17","description":"","hostname":"hostname","hypervisor_id":null,"id":7,"mac":"52:54:00:6c:6a:ac","memory":2048,"name":"Ubuntu-16-04","os_name":"ubuntu","os_version":"16.04","private_ipv4":"","status":"running","storage":[{"created_at":"2018-02-22 11:10:37.793878","id":5,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a"}],"tags":null,"updated_at":"2018-02-22 13:48:21","user_id":8,"username":"example","uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a","vcpu":2}`,
			StatusCode: http.StatusOK,
			Error:      DiskNotFoundError(),
		},
		{
			RequestData: map[string]interface{}{
				"vm_uuid": "valid-vm-uuid",
				"uuid":    "this-is-some-valid-disk-uuid",
			},
			Body:       `{"backup":true,"billing_account":6,"created_at":"2018-02-22 11:10:17","description":"","hostname":"hostname","hypervisor_id":null,"id":7,"mac":"52:54:00:6c:6a:ac","memory":2048,"name":"Ubuntu-16-04","os_name":"ubuntu","os_version":"16.04","private_ipv4":"","status":"running","storage":[{"created_at":"2018-02-22 11:10:37.793878","id":5,"name":"sda","pool":"default2","primary":true,"replica":[],"shared":false,"size":20,"type":"block","updated_at":null,"user_id":8,"uuid":"this-is-some-valid-disk-uuid"}],"tags":null,"updated_at":"2018-02-22 13:48:21","user_id":8,"username":"example","uuid":"f80b1d62-ffe4-43ef-9210-60f05445456a","vcpu":2}`,
			StatusCode: http.StatusOK,
			Error:      DiskNotFoundError(),
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}
		testDiskApi.Bind(test.RequestData["vm_uuid"].(string))
		err := testDiskApi.Get(test.RequestData["uuid"].(string))
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}

func TestCreateDisk(t *testing.T) {
	testDiskApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"vm_uuid":   "valid-vm-uuid",
				"disk_size": 50,
			},
			Body:       `{"created_at":"2019-08-14 13:57:44","name":"vdc","pool":"default","primary":false,"replica":[],"shared":false,"size":50,"type":"block","uuid":"3d91aa31-16ec-44ee-b8b3-22a0bda6559e"}`,
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
		testDiskApi.Bind(test.RequestData["vm_uuid"].(string))
		err := testDiskApi.Create(test.RequestData["disk_size"].(int))
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}
