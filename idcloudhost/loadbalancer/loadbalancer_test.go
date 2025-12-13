package loadbalancer

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
	mockHttpClient      = &HTTPClientMock{}
	testLoadBalancerAPI = LoadBalancerAPI{}
	loc                 = "jkt01"
)

func TestListLoadBalancers(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		Body       string
		StatusCode int
		Error      error
	}{
		{
			Body:       `[{"uuid":"72fa225f-5cb1-45cf-83ed-b0f441b0b815","display_name":"xyz1234","user_id":4384,"billing_account_id":1200132376,"created_at":"2025-12-13 11:54:16","updated_at":"2025-12-13 11:54:16","is_deleted":false,"deleted_at":null,"private_address":"10.4.207.254","network_uuid":"730bf645-9b36-44b8-8ca1-46d2480cc0d6","forwarding_rules":[{"protocol":"TCP","uuid":"a96ee226-7204-47b4-8782-02fbbdb6d4e7","created_at":"2025-12-13 11:54:16","source_port":80,"target_port":80,"settings":{"connection_limit":10000,"session_persistence":"SOURCE_IP"}}],"targets":[{"created_at":"2025-12-13 11:54:16","target_uuid":"70517643-b046-48e3-9bae-83dc2c143beb","target_type":"vm","target_ip_address":"10.4.207.133"}]}]`,
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
		err := testLoadBalancerAPI.ListLoadBalancers()
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancerList)
			assert.Equal(t, "72fa225f-5cb1-45cf-83ed-b0f441b0b815", testLoadBalancerAPI.LoadBalancerList[0].UUID)
			assert.Equal(t, "xyz1234", testLoadBalancerAPI.LoadBalancerList[0].DisplayName)
			assert.Equal(t, "10.4.207.254", testLoadBalancerAPI.LoadBalancerList[0].PrivateAddress)
			assert.Equal(t, "730bf645-9b36-44b8-8ca1-46d2480cc0d6", testLoadBalancerAPI.LoadBalancerList[0].NetworkUUID)
		}
	}
}

func TestGetLoadBalancer(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "72fa225f-5cb1-45cf-83ed-b0f441b0b815",
			},
			Body:       `{"uuid":"72fa225f-5cb1-45cf-83ed-b0f441b0b815","display_name":"xyz1234","user_id":4384,"billing_account_id":1200132376,"created_at":"2025-12-13 11:54:16","updated_at":"2025-12-13 11:54:16","is_deleted":false,"deleted_at":null,"private_address":"10.4.207.254","network_uuid":"730bf645-9b36-44b8-8ca1-46d2480cc0d6","forwarding_rules":[{"protocol":"TCP","uuid":"a96ee226-7204-47b4-8782-02fbbdb6d4e7","created_at":"2025-12-13 11:54:16","source_port":80,"target_port":80,"settings":{"connection_limit":10000,"session_persistence":"SOURCE_IP"}}],"targets":[{"created_at":"2025-12-13 11:54:16","target_uuid":"70517643-b046-48e3-9bae-83dc2c143beb","target_type":"vm","target_ip_address":"10.4.207.133"}]}`,
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
		uuid := test.RequestData["uuid"].(string)
		err := testLoadBalancerAPI.GetLoadBalancer(uuid)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.Equal(t, uuid, testLoadBalancerAPI.LoadBalancer.UUID)
			assert.Equal(t, "xyz1234", testLoadBalancerAPI.LoadBalancer.DisplayName)
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.Targets)
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.ForwardingRules)
			assert.Equal(t, "10.4.207.254", testLoadBalancerAPI.LoadBalancer.PrivateAddress)
			assert.Equal(t, "vm", testLoadBalancerAPI.LoadBalancer.Targets[0].TargetType)
			assert.Equal(t, "70517643-b046-48e3-9bae-83dc2c143beb", testLoadBalancerAPI.LoadBalancer.Targets[0].TargetUUID)
		}
	}
}

func TestCreateLoadBalancer(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData *CreateLoadBalancerRequest
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: &CreateLoadBalancerRequest{
				ReservePublicIP: true,
				NetworkUUID:     "730bf645-9b36-44b8-8ca1-46d2480cc0d6",
				Targets: []CreateTargetRequest{
					{
						TargetType: "vm",
						TargetUUID: "70517643-b046-48e3-9bae-83dc2c143beb",
					},
				},
				Rules: []CreateRuleRequest{
					{
						SourcePort: 80,
						TargetPort: 80,
					},
				},
				DisplayName:      "xyz1234",
				BillingAccountID: 1200132376,
			},
			Body:       `{"uuid":"72fa225f-5cb1-45cf-83ed-b0f441b0b815","display_name":"xyz1234","user_id":4384,"billing_account_id":1200132376,"created_at":"2025-12-13 11:54:16","updated_at":"2025-12-13 11:54:16","is_deleted":false,"deleted_at":null,"private_address":"10.4.207.254","network_uuid":"730bf645-9b36-44b8-8ca1-46d2480cc0d6","forwarding_rules":[{"protocol":"TCP","uuid":"a96ee226-7204-47b4-8782-02fbbdb6d4e7","created_at":"2025-12-13 11:54:16","source_port":80,"target_port":80,"settings":{"connection_limit":10000,"session_persistence":"SOURCE_IP"}}],"targets":[{"created_at":"2025-12-13 11:54:16","target_uuid":"70517643-b046-48e3-9bae-83dc2c143beb","target_type":"vm","target_ip_address":"10.4.207.133"}]}`,
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
		err := testLoadBalancerAPI.CreateLoadBalancer(test.RequestData)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.Equal(t, test.RequestData.DisplayName, testLoadBalancerAPI.LoadBalancer.DisplayName)
			assert.Equal(t, test.RequestData.BillingAccountID, testLoadBalancerAPI.LoadBalancer.BillingAccountID)
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.UUID)
			assert.Equal(t, test.RequestData.NetworkUUID, testLoadBalancerAPI.LoadBalancer.NetworkUUID)
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.Targets)
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.ForwardingRules)
		}
	}
}

func TestModifyLoadBalancer(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid":         "72fa225f-5cb1-45cf-83ed-b0f441b0b815",
				"display_name": "renamed-lb",
			},
			Body:       `{"uuid":"72fa225f-5cb1-45cf-83ed-b0f441b0b815","display_name":"renamed-lb","user_id":4384,"billing_account_id":1200132376,"created_at":"2025-12-13 11:54:16","updated_at":"2025-12-13 12:00:00","is_deleted":false,"deleted_at":null,"private_address":"10.4.207.254","network_uuid":"730bf645-9b36-44b8-8ca1-46d2480cc0d6","forwarding_rules":[],"targets":[]}`,
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
		uuid := test.RequestData["uuid"].(string)
		displayName := test.RequestData["display_name"].(string)
		err := testLoadBalancerAPI.RenameLoadBalancer(uuid, displayName)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.Equal(t, displayName, testLoadBalancerAPI.LoadBalancer.DisplayName)
			assert.Equal(t, uuid, testLoadBalancerAPI.LoadBalancer.UUID)
		}
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "72fa225f-5cb1-45cf-83ed-b0f441b0b815",
			},
			Body:       ``,
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
		uuid := test.RequestData["uuid"].(string)
		err := testLoadBalancerAPI.DeleteLoadBalancer(uuid)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
	}
}

func TestAddTarget(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "72fa225f-5cb1-45cf-83ed-b0f441b0b815",
				"target": &CreateTargetRequest{
					TargetType: "vm",
					TargetUUID: "70517643-b046-48e3-9bae-83dc2c143beb",
				},
			},
			Body:       `{"created_at":"2025-12-13 11:55:09","target_uuid":"70517643-b046-48e3-9bae-83dc2c143beb","target_type":"vm","target_ip_address":null}`,
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
		uuid := test.RequestData["uuid"].(string)
		target := test.RequestData["target"].(*CreateTargetRequest)
		err := testLoadBalancerAPI.AddTarget(uuid, target)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotNil(t, testLoadBalancerAPI.Target)
			assert.Equal(t, target.TargetUUID, testLoadBalancerAPI.Target.TargetUUID)
			assert.Equal(t, target.TargetType, testLoadBalancerAPI.Target.TargetType)
		}
	}
}

func TestRemoveTarget(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid":        "72fa225f-5cb1-45cf-83ed-b0f441b0b815",
				"target_uuid": "70517643-b046-48e3-9bae-83dc2c143beb",
			},
			Body:       ``,
			StatusCode: http.StatusNoContent,
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
		uuid := test.RequestData["uuid"].(string)
		targetUUID := test.RequestData["target_uuid"].(string)
		err := testLoadBalancerAPI.RemoveTarget(uuid, targetUUID)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
	}
}

func TestAddRule(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "72fa225f-5cb1-45cf-83ed-b0f441b0b815",
				"rule": &CreateRuleRequest{
					SourcePort: 443,
					TargetPort: 8443,
				},
			},
			Body:       `{"protocol":"TCP","uuid":"5b0d63c4-6998-49c8-a09d-651f5763a599","created_at":"2025-12-13 11:56:16","source_port":443,"target_port":8443,"settings":{"connection_limit":10000,"session_persistence":"SOURCE_IP"}}`,
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
		uuid := test.RequestData["uuid"].(string)
		rule := test.RequestData["rule"].(*CreateRuleRequest)
		err := testLoadBalancerAPI.AddRule(uuid, rule)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotNil(t, testLoadBalancerAPI.ForwardingRule)
			assert.Equal(t, rule.SourcePort, testLoadBalancerAPI.ForwardingRule.SourcePort)
			assert.Equal(t, rule.TargetPort, testLoadBalancerAPI.ForwardingRule.TargetPort)
			assert.Equal(t, "TCP", testLoadBalancerAPI.ForwardingRule.Protocol)
			assert.NotEmpty(t, testLoadBalancerAPI.ForwardingRule.UUID)
		}
	}
}

func TestRemoveRule(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid":      "72fa225f-5cb1-45cf-83ed-b0f441b0b815",
				"rule_uuid": "a96ee226-7204-47b4-8782-02fbbdb6d4e7",
			},
			Body:       ``,
			StatusCode: http.StatusNoContent,
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
		uuid := test.RequestData["uuid"].(string)
		ruleUUID := test.RequestData["rule_uuid"].(string)
		err := testLoadBalancerAPI.RemoveRule(uuid, ruleUUID)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
	}
}

func TestChangeBillingAccount(t *testing.T) {
	if err := testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize loadbalancer api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid":               "72fa225f-5cb1-45cf-83ed-b0f441b0b815",
				"billing_account_id": 5678,
			},
			Body:       `{"uuid":"72fa225f-5cb1-45cf-83ed-b0f441b0b815","display_name":"xyz1234","user_id":4384,"billing_account_id":5678,"created_at":"2025-12-13 11:54:16","updated_at":"2025-12-13 12:00:00","is_deleted":false,"deleted_at":null,"private_address":"10.4.207.254","network_uuid":"730bf645-9b36-44b8-8ca1-46d2480cc0d6","forwarding_rules":[],"targets":[]}`,
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
		uuid := test.RequestData["uuid"].(string)
		billingAccountID := test.RequestData["billing_account_id"].(int)
		err := testLoadBalancerAPI.ChangeBillingAccount(uuid, billingAccountID)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.Equal(t, billingAccountID, testLoadBalancerAPI.LoadBalancer.BillingAccountID)
		}
	}
}
