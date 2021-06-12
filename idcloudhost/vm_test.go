package idcloudhost

import (
	"encoding/json"
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

const userAuthToken = "xxxxx"

var (
	c          = &HTTPClientMock{}
	v          = VirtualMachineAPI{}
	targetUuid = "validuuid-cdb2-1234-b6f8-8f7deadbeef0"
	loc        = "jkt01"
)

func TestGetVMbyUUID(t *testing.T) {
	v.Init(c, userAuthToken, loc)
	if err := v.Get(targetUuid); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(v.VMMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(b))
}
func TestModify(t *testing.T) {
	v.Init(c, userAuthToken, loc)
	propertyMap := map[string]interface{}{
		"uuid": targetUuid, "vcpu": 1, "ram": 1536, "name": "testvm",
	}
	if err := v.Modify(propertyMap); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(v.VMMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(b))
}

func TestListAllVMs(t *testing.T) {
	v.Init(c, userAuthToken, loc)
	if err := v.ListAll(); err != nil {
		t.Fatal(err)
	}
	var inInterface []map[string]interface{}
	inrec, _ := json.Marshal(v.VMListMap)
	log.Println(json.Unmarshal(inrec, &inInterface))

	// iterate through inrecs
	for field, val := range inInterface {
		log.Println("KV Pair: ", field, val)
	}
}
func TestCreateVM(t *testing.T) {
	v.Init(c, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"backup":          false,
				"name":            "testvm",
				"os_name":         "ubuntu",
				"os_version":      "16.04",
				"disks":           20,
				"vcpu":            1,
				"ram":             1024,
				"username":        "example",
				"password":        "Password123",
				"billing_account": 9999,
			},
			Body:       `{"backup": false, "name": "testvm", "os_name": "ubuntu", "os_version": "16.04", "disks": [], "vcpu": 1, "ram": 1024, "username": "example", "password": "Password123", "billing_account": 9999, }`,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
		{
			RequestData: map[string]interface{}{
				"name": "incomplete-vm",
			},
			Body:       ``,
			StatusCode: http.StatusBadRequest,
			Error:      BadRequestError(),
		},
	}
	for _, test := range testCases {
		c.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := v.Create(test.RequestData)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
	}
}
