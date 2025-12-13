//go:build !integration

package floatingip

import (
	"io"
	"net/http"
	"strings"
	"testing"

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
	mockHttpClient    = &HTTPClientMock{}
	testFloatingIPAPI = FloatingIPAPI{}
	loc               = "jkt01"
)

func TestGetIP(t *testing.T) {
	if err := testFloatingIPAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize floating ip api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"ip_address": "1.1.1.1",
			},
			Body:       `{"id":1,"address":"1.1.1.1","user_id":666,"billing_account_id":1,"type":"public","network_id":null,"name":"Test IP","enabled":true,"created_at":"2019-10-31 10:52:19","updated_at":"2019-11-01 10:22:19","assigned_to":"88e5a11b-9c89-4986-99c7-90d43499317c"}`,
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
		ipAddress := test.RequestData["ip_address"].(string)
		err := testFloatingIPAPI.Get(ipAddress)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		assert.Equal(t, testFloatingIPAPI.FloatingIP.Address, ipAddress)
	}
}

func TestCreateFloatingIP(t *testing.T) {
	if err := testFloatingIPAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize floating ip api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"name":               "Test Create IP",
				"billing_account_id": 1111,
			},
			Body:       `{"id":1,"address":"1.1.1.1","user_id":666,"billing_account_id":1111,"type":"public","network_id":null,"name":"Test Create IP","enabled":true,"created_at":"2019-10-31 10:52:19","updated_at":"2019-11-01 10:22:19","assigned_to":"88e5a11b-9c89-4986-99c7-90d43499317c"}`,
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
		ipName := test.RequestData["name"].(string)
		billingAcc := test.RequestData["billing_account_id"].(int)
		err := testFloatingIPAPI.Get(ipName)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		assert.Equal(t, testFloatingIPAPI.FloatingIP.Name, ipName)
		assert.Equal(t, testFloatingIPAPI.FloatingIP.BillingAccount, billingAcc)
	}
}
