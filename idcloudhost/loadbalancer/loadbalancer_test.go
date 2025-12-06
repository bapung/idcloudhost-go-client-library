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
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		Body       string
		StatusCode int
		Error      error
	}{
		{
			Body:       `[{"id":1,"uuid":"lb-uuid-123","name":"test-lb","billing_account_id":1234,"user_id":1,"target_ips":[],"forward_rules":[{"id":1,"name":"http-rule","protocol":"http","frontend_port":80,"backend_port":8080,"health_check":true,"health_timeout":30}],"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}]`,
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
			assert.Equal(t, "lb-uuid-123", testLoadBalancerAPI.LoadBalancerList[0].UUID)
			assert.Equal(t, "test-lb", testLoadBalancerAPI.LoadBalancerList[0].Name)
		}
	}
}

func TestGetLoadBalancer(t *testing.T) {
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "lb-uuid-123",
			},
			Body:       `{"id":1,"uuid":"lb-uuid-123","name":"test-lb","billing_account_id":1234,"user_id":1,"target_ips":[{"vm_id":1,"vm_uuid":"vm-uuid-1","ip_addr":"10.0.0.1","name":"vm1"}],"forward_rules":[{"id":1,"name":"http-rule","protocol":"http","frontend_port":80,"backend_port":8080,"health_check":true,"health_timeout":30}],"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}`,
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
			assert.Equal(t, "test-lb", testLoadBalancerAPI.LoadBalancer.Name)
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.TargetIPs)
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.ForwardRules)
		}
	}
}

func TestCreateLoadBalancer(t *testing.T) {
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData *LoadBalancer
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: &LoadBalancer{
				Name:           "new-lb",
				BillingAccount: 1234,
				ForwardRules: []ForwardingRule{
					{
						Name:         "http-rule",
						Protocol:     "http",
						FrontendPort: 80,
						BackendPort:  8080,
						HealthCheck:  true,
					},
				},
			},
			Body:       `{"id":1,"uuid":"lb-uuid-new","name":"new-lb","billing_account_id":1234,"user_id":1,"target_ips":[],"forward_rules":[{"id":1,"name":"http-rule","protocol":"http","frontend_port":80,"backend_port":8080,"health_check":true,"health_timeout":30}],"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}`,
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
			assert.Equal(t, test.RequestData.Name, testLoadBalancerAPI.LoadBalancer.Name)
			assert.Equal(t, test.RequestData.BillingAccount, testLoadBalancerAPI.LoadBalancer.BillingAccount)
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.UUID)
		}
	}
}

func TestRenameLoadBalancer(t *testing.T) {
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "lb-uuid-123",
				"name": "renamed-lb",
			},
			Body:       `{"id":1,"uuid":"lb-uuid-123","name":"renamed-lb","billing_account_id":1234,"user_id":1,"target_ips":[],"forward_rules":[],"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 11:00:00"}`,
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
		name := test.RequestData["name"].(string)
		err := testLoadBalancerAPI.RenameLoadBalancer(uuid, name)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.Equal(t, name, testLoadBalancerAPI.LoadBalancer.Name)
			assert.Equal(t, uuid, testLoadBalancerAPI.LoadBalancer.UUID)
		}
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "lb-uuid-123",
			},
			Body:       `{"success":true}`,
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
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "lb-uuid-123",
				"target": &Target{
					VMUUID: "vm-uuid-1",
				},
			},
			Body:       `{"id":1,"uuid":"lb-uuid-123","name":"test-lb","billing_account_id":1234,"user_id":1,"target_ips":[{"vm_id":1,"vm_uuid":"vm-uuid-1","ip_addr":"10.0.0.1","name":"vm1"}],"forward_rules":[],"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}`,
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
		target := test.RequestData["target"].(*Target)
		err := testLoadBalancerAPI.AddTarget(uuid, target)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.TargetIPs)
			assert.Equal(t, target.VMUUID, testLoadBalancerAPI.LoadBalancer.TargetIPs[0].VMUUID)
		}
	}
}

func TestRemoveTarget(t *testing.T) {
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid":    "lb-uuid-123",
				"vm_uuid": "vm-uuid-1",
			},
			Body:       `{"success":true}`,
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
		vmUUID := test.RequestData["vm_uuid"].(string)
		err := testLoadBalancerAPI.RemoveTarget(uuid, vmUUID)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
	}
}

func TestAddRule(t *testing.T) {
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "lb-uuid-123",
				"rule": &ForwardingRule{
					Name:         "https-rule",
					Protocol:     "https",
					FrontendPort: 443,
					BackendPort:  8443,
					HealthCheck:  true,
				},
			},
			Body:       `{"id":1,"uuid":"lb-uuid-123","name":"test-lb","billing_account_id":1234,"user_id":1,"target_ips":[],"forward_rules":[{"id":1,"name":"https-rule","protocol":"https","frontend_port":443,"backend_port":8443,"health_check":true,"health_timeout":30}],"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}`,
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
		rule := test.RequestData["rule"].(*ForwardingRule)
		err := testLoadBalancerAPI.AddRule(uuid, rule)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotEmpty(t, testLoadBalancerAPI.LoadBalancer.ForwardRules)
			assert.Equal(t, rule.Name, testLoadBalancerAPI.LoadBalancer.ForwardRules[0].Name)
			assert.Equal(t, rule.Protocol, testLoadBalancerAPI.LoadBalancer.ForwardRules[0].Protocol)
		}
	}
}

func TestRemoveRule(t *testing.T) {
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid":    "lb-uuid-123",
				"rule_id": 1,
			},
			Body:       `{"success":true}`,
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
		ruleID := test.RequestData["rule_id"].(int)
		err := testLoadBalancerAPI.RemoveRule(uuid, ruleID)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
	}
}

func TestChangeBillingAccount(t *testing.T) {
	testLoadBalancerAPI.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid":               "lb-uuid-123",
				"billing_account_id": 5678,
			},
			Body:       `{"id":1,"uuid":"lb-uuid-123","name":"test-lb","billing_account_id":5678,"user_id":1,"target_ips":[],"forward_rules":[],"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 11:00:00"}`,
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
			assert.Equal(t, billingAccountID, testLoadBalancerAPI.LoadBalancer.BillingAccount)
		}
	}
}
