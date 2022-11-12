package disk

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
)

type HTTPClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (H HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	return H.DoFunc(r)
}

const userAuthToken = "xxxxx"

var (
	mockHttpClient = &HTTPClientMock{}
	loc            = "jkt01"
	testDiskApi    = DiskAPI{}
)

func TestGetDisk(t *testing.T) {
	testDiskApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		vmUUID		string
		UUID 		string
		DiskList  	[]DiskStorage
		Error		error
	}{
		{
			vmUUID: "someVMuuid",
			UUID: "okeuuid",
			DiskList: []DiskStorage{
				{
					"aaa", 1, "aaaa", "aaa", true, []string{}, false, 20, "aaa", "aaa", 123, "okeuuid",
				},
				{
					"aaa", 2, "aaaa", "aaa", true, []string{}, false, 20, "aaa", "aaa", 123, "falseuuid",
				},
			},
			Error: nil,
		},
	}
	for _, test := range testCases {
		testDiskApi.Bind(test.vmUUID)
		err := testDiskApi.Get(test.UUID, &test.DiskList)
		if err != nil && err != test.Error {
			t.Fatalf("want %v, got %v", err, test.Error)
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
		diskSize := test.RequestData["disk_size"].(int)
		err := testDiskApi.Create(diskSize)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		assert.Equal(t, testDiskApi.Disk.SizeGB, diskSize)
	}
}

func TestModifyDisk(t *testing.T) {
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
				"disk_uuid": "valid-disk-uuid",
				"disk_size": 50,
			},
			Body:       `{"created_at":"2019-08-14 13:57:44","name":"vdc","pool":"default","primary":false,"replica":[],"shared":false,"size":50,"type":"block","uuid":"valid-disk-uuid"}`,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
		{
			RequestData: map[string]interface{}{
				"vm_uuid":   "valid-vm-uuid",
				"disk_uuid": "non-exist-disk-uuid",
				"disk_size": 50,
			},
			Body:       `{"created_at":"2019-08-14 13:57:44","name":"vdc","pool":"default","primary":false,"replica":[],"shared":false,"size":50,"type":"block","uuid":"valid-disk-uuid"}`,
			StatusCode: http.StatusNotFound,
			Error:      errors.New(fmt.Sprintf("%v",http.StatusNotFound)),
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
		diskUUID := test.RequestData["disk_uuid"].(string)
		diskSize := test.RequestData["disk_size"].(int)
		err := testDiskApi.Modify(diskUUID, diskSize)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		if test.StatusCode == http.StatusOK {
			assert.Equal(t, testDiskApi.Disk.UUID, diskUUID)
			assert.Equal(t, testDiskApi.Disk.SizeGB, diskSize)
		}
	}
}

func TestDeleteDisk(t *testing.T) {
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
				"disk_uuid": "valid-disk-uuid",
			},
			Body:       `{ "success": true }`,
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
		err := testDiskApi.Delete(test.RequestData["disk_uuid"].(string))
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}
