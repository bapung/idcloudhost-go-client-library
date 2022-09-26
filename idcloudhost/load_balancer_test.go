package idcloudhost

import (
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	//"github.com/stretchr/testify/assert"
)

var testLoadBalancerApi LoadBalancerAPI

func TestListAllLB(t *testing.T) {
	testLoadBalancerApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		Body       string
		StatusCode int
		Error      error
	}{
		{
			Body:       `[{"uuid":"438ac62f-e97b-4ef0-8940-507b9e94af43","network_uuid":"438ac62f-e97b-4ef0-8940-507b9e94af43","user_id":268,"billing_account_id":130157,"created_at":"2022-07-12 14:21:06","updated_at":"2022-07-12 14:21:06","is_deleted":false,"private_address":"10.112.231.192","forwarding_rules":[{"uuid":"b3f28feb-c91e-4601-a6b6-267fa98dc121","protocol":"TCP","created_at":"2022-07-12 14:21:06","source_port":8080,"target_port":8080,"settings":{"connection_limit":10000,"session_persistence":"SOURCE_IP"}}],"targets":[{"created_at":"2022-07-12 14:21:06","target_uuid":"145cc106-e067-419a-85fd-333ded30f169","target_type":"vm","target_ip_address":"10.61.10.2"}]}]`,
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
		err := testLoadBalancerApi.ListAll(true)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}

func TestCreateLBB(t *testing.T) {
	testLoadBalancerApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData LoadBalancer
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: LoadBalancer{
				DisplayName:     "my LB",
				BillingAccount:  130157,
				NetworkUUID:     "438ac62f-e97b-4ef0-8940-507b9e94af43",
				ReservePublicIP: true,
				ForwardingRules: []ForwardingRule{
					{
						SourcePort: 8080,
						TargetPort: 80,
					},
				},
				Targets: []ForwardingTarget{
					{
						TargetUUID: "145cc106-e067-419a-85fd-333ded30f169",
						TargetType: "vm",
					},
					ForwardingTarget{
						TargetUUID: "e9717243-59df-4847-bd50-dca5b090432b",
						TargetType: "vm",
					},
				},
			},
			Body:       `{"uuid":"438ac62f-e97b-4ef0-8940-507b9e94af43","display_name":"my LB","network_uuid":"438ac62f-e97b-4ef0-8940-507b9e94af43","user_id":268,"billing_account_id":130157,"created_at":"2022-07-12 14:21:06","updated_at":"2022-07-12 14:21:06","is_deleted":false,"private_address":"10.112.231.192","forwarding_rules":[{"uuid":"b3f28feb-c91e-4601-a6b6-267fa98dc121","protocol":"TCP","created_at":"2022-07-12 14:21:06","source_port":8080,"target_port":8080,"settings":{"connection_limit":10000,"session_persistence":"SOURCE_IP"}}],"targets":[{"created_at":"2022-07-12 14:21:06","target_uuid":"145cc106-e067-419a-85fd-333ded30f169","target_type":"vm","target_ip_address":"10.61.10.2"}]}`,
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

		err := testLoadBalancerApi.Create(true, &test.RequestData)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
		log.Println(testLoadBalancerApi.LoadBalancer)
	}
}

func TestDeleteLB(t *testing.T) {
	testLoadBalancerApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		LBUUID     string
		StatusCode int
		Error      error
	}{
		{
			LBUUID:     "438ac62f-e97b-4ef0-8940-507b9e94af43",
			StatusCode: http.StatusOK,
			Error:      nil,
		},
	}

	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: test.StatusCode,
				Body:       nil,
			}, nil
		}
		err := testLoadBalancerApi.Delete(test.LBUUID)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}

func TestAddForwardingTarget(t *testing.T) {
	testLoadBalancerApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]string
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]string{
				"LBUUID":     "438ac62f-e97b-4ef0-8940-507b9e94af43",
				"TargetUUID": "145cc106-e067-419a-85fd-333ded30f169",
				"TargetType": "vm",
			},
			Body:       `{"created_at":"2022-07-12 14:21:06","target_uuid":"145cc106-e067-419a-85fd-333ded30f169","target_type":"vm","target_ip_address":"10.61.10.2"}`,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
	}

	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: test.StatusCode,
				Body:       io.NopCloser(strings.NewReader(test.Body)),
			}, nil
		}
		err := testLoadBalancerApi.AddForwardingTarget(
			test.RequestData["LBUUID"],
			test.RequestData["TargetUUID"],
			test.RequestData["TargetType"],
		)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}

func TestDeleteForwardingRule(t *testing.T) {
	testLoadBalancerApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]string
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]string{
				"LBUUID":   "438ac62f-e97b-4ef0-8940-507b9e94af43",
				"RuleUUID": "145cc106-e067-419a-85fd-333ded30f169",
			},
			StatusCode: http.StatusOK,
			Error:      nil,
		},
	}

	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: test.StatusCode,
			}, nil
		}
		err := testLoadBalancerApi.DeleteForwardingRule(
			test.RequestData["LBUUID"],
			test.RequestData["RuleUUID"],
		)
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}
