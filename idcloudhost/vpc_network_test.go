package idcloudhost

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

var testNetworkAPI VPCNetworkAPI

func TestListNetwork(t *testing.T) {
	testNetworkAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{},
			Body:        `[{"vlan_id":1175,"subnet":"10.5.155.0/24","name":"My Network","created_at":"2021-05-26 15:14:52","updated_at":"2021-08-02 13:03:40","uuid":"cf99ed55-608f-438a-a371-32a2f5813cbc","type":"private","is_default":true,"vm_uuids":["56accb3c-2b45-45b0-af83-90afef6974e8"],"resources_count":1}]`,
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
		err := testNetworkAPI.List()
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}

func TestGetNetwork(t *testing.T) {
	testNetworkAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		UUID       string
		Body       string
		StatusCode int
		Error      error
	}{
		{
			UUID:       "cf99ed55-608f-438a-a371-32a2f5813cbc",
			Body:       `{"vlan_id":1175,"subnet":"10.5.155.0/24","name":"My Network","created_at":"2021-05-26 15:14:52","updated_at":"2021-08-02 13:03:40","uuid":"cf99ed55-608f-438a-a371-32a2f5813cbc","type":"private","is_default":true,"vm_uuids":["56accb3c-2b45-45b0-af83-90afef6974e8"],"resources_count":1}`,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
		{
			UUID:       "non-valid-uuid",
			Body:       `{"timestamp":1627912221481,"status":500,"error":"Internal Server Error","message":"Network uuid is invalid.","path":"/v1/network/non-valid-uuid"}`,
			StatusCode: http.StatusInternalServerError,
			Error:      UnknownError(),
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := testNetworkAPI.Get(test.UUID)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}
func TestCreateNetwork(t *testing.T) {
	testNetworkAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		Name       string
		Body       string
		StatusCode int
		Error      error
	}{
		{
			Name:       "test-network-created",
			Body:       `{"vlan_id":1175,"subnet":"10.5.155.0/24","name":"test-network-created","created_at":"2021-05-26 15:14:52","updated_at":"2021-08-02 13:03:40","uuid":"cf99ed55-608f-438a-a371-32a2f5813cbc","type":"private","is_default":true,"vm_uuids":[],"resources_count":1}`,
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

		err := testNetworkAPI.Create(test.Name)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}

func TestDeleteNetwork(t *testing.T) {
	testNetworkAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		UUID       string
		Body       string
		StatusCode int
		Error      error
	}{
		{
			UUID:       "test-non-default-no-resource-network",
			Body:       ``,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
		{
			UUID:       "test-non-exist-network",
			Body:       `{"timestamp":1628085172112,"status":500,"error":"Internal Server Error","message":"Network uuid is invalid.","path":"/v1/network/82bf2dd0-c9ab-429a-a6a5-8e17f56caeaf"}`,
			StatusCode: http.StatusInternalServerError,
			Error:      UnknownError(),
		},
		{
			UUID:       "test-default-network",
			Body:       `{"timestamp":1628085453108,"status":500,"error":"Internal Server Error","message":"Default network cannot be deleted.","path":"/v1/network/cf99ed55-608f-438a-a371-32a2f5813cbc"}`,
			StatusCode: http.StatusInternalServerError,
			Error:      UnknownError(),
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := testNetworkAPI.Delete(test.UUID)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}
