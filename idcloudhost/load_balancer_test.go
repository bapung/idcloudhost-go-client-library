package idcloudhost

import (
	"testing"
	"net/http"
	"io"
	"strings"
	//"github.com/stretchr/testify/assert"

)

var testLoadBalancerApi LoadBalancerAPI

func TestListAllLB(t *testing.T) {
	testLoadBalancerApi.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		Body		string
		StatusCode	int
		Error		error
	}{
		{
			Body: 		`[{"uuid":"438ac62f-e97b-4ef0-8940-507b9e94af43","network_uuid":"438ac62f-e97b-4ef0-8940-507b9e94af43","user_id":268,"billing_account_id":130157,"created_at":"2022-07-12 14:21:06","updated_at":"2022-07-12 14:21:06","is_deleted":false,"private_address":"10.112.231.192","forwarding_rules":[{"uuid":"b3f28feb-c91e-4601-a6b6-267fa98dc121","protocol":"TCP","created_at":"2022-07-12 14:21:06","source_port":8080,"target_port":8080,"settings":{"connection_limit":10000,"session_persistence":"SOURCE_IP"}}],"targets":[{"created_at":"2022-07-12 14:21:06","target_uuid":"145cc106-e067-419a-85fd-333ded30f169","target_type":"vm","target_ip_address":"10.61.10.2"}]}]`,
			StatusCode: http.StatusOK,
			Error:		nil,
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