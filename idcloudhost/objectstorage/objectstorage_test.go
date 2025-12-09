package objectstorage

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
	loc = "jkt01"
)

func setupMockClient(responseBody string) *HTTPClientMock {
	return &HTTPClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
			}, nil
		},
	}
}

func TestGetS3Info(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		Body       string
		StatusCode int
		Error      error
	}{
		{
			Body:       `{"id":1,"user_id":123,"billing_account_id":1234,"total_storage_gb":100.0,"used_storage_gb":25.5,"storage_endpoint":"https://s3.example.com","service_endpoint":"https://service.example.com","created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}`,
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
		err := testObjectStorageAPI.GetS3Info()
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotNil(t, testObjectStorageAPI.S3User)
			assert.Equal(t, 1, testObjectStorageAPI.S3User.ID)
			assert.Equal(t, 123, testObjectStorageAPI.S3User.UserID)
			assert.Equal(t, 100.0, testObjectStorageAPI.S3User.TotalStorageGB)
			assert.Equal(t, 25.5, testObjectStorageAPI.S3User.UsedStorageGB)
		}
	}
}

func TestCreateBucket(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"name":               "test-bucket",
				"acl":                "private",
				"billing_account_id": 1234,
			},
			Body:       `{"id":1,"name":"test-bucket","billing_account_id":1234,"user_id":123,"size_bytes":0,"region":"jkt01","acl":"private","created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}`,
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
		name := test.RequestData["name"].(string)
		acl := test.RequestData["acl"].(string)
		billingAccountID := test.RequestData["billing_account_id"].(int)
		err := testObjectStorageAPI.CreateBucket(name, acl, billingAccountID)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.Equal(t, name, testObjectStorageAPI.Bucket.Name)
			assert.Equal(t, acl, testObjectStorageAPI.Bucket.ACL)
			assert.Equal(t, billingAccountID, testObjectStorageAPI.Bucket.BillingAccount)
		}
	}
}

func TestModifyBucket(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"name": "test-bucket",
				"acl":  "public-read",
			},
			Body:       `{"id":1,"name":"test-bucket","billing_account_id":1234,"user_id":123,"size_bytes":1024,"region":"jkt01","acl":"public-read","created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 11:00:00"}`,
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
		name := test.RequestData["name"].(string)
		acl := test.RequestData["acl"].(string)
		err := testObjectStorageAPI.ModifyBucket(name, acl)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.Equal(t, acl, testObjectStorageAPI.Bucket.ACL)
			assert.Equal(t, name, testObjectStorageAPI.Bucket.Name)
		}
	}
}

func TestDeleteBucket(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"name": "test-bucket",
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
		name := test.RequestData["name"].(string)
		err := testObjectStorageAPI.DeleteBucket(name)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
	}
}

func TestGetBucket(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"name": "test-bucket",
			},
			Body:       `{"id":1,"name":"test-bucket","billing_account_id":1234,"user_id":123,"size_bytes":2048,"region":"jkt01","acl":"private","created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}`,
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
		name := test.RequestData["name"].(string)
		err := testObjectStorageAPI.GetBucket(name)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.Equal(t, name, testObjectStorageAPI.Bucket.Name)
			assert.Equal(t, int64(2048), testObjectStorageAPI.Bucket.SizeBytes)
		}
	}
}

func TestListBuckets(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		Body       string
		StatusCode int
		Error      error
	}{
		{
			Body:       `[{"id":1,"name":"bucket-1","billing_account_id":1234,"user_id":123,"size_bytes":1024,"region":"jkt01","acl":"private","created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"},{"id":2,"name":"bucket-2","billing_account_id":1234,"user_id":123,"size_bytes":2048,"region":"jkt01","acl":"public-read","created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}]`,
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
		err := testObjectStorageAPI.ListBuckets()
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotEmpty(t, testObjectStorageAPI.Buckets)
			assert.Equal(t, 2, len(testObjectStorageAPI.Buckets))
			assert.Equal(t, "bucket-1", testObjectStorageAPI.Buckets[0].Name)
			assert.Equal(t, "bucket-2", testObjectStorageAPI.Buckets[1].Name)
		}
	}
}

func TestGetKeys(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		Body       string
		StatusCode int
		Error      error
	}{
		{
			Body:       `[{"access_key":"AKIAIOSFODNN7EXAMPLE","secret_key":"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY","user_id":123,"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}]`,
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
		err := testObjectStorageAPI.GetKeys()
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotEmpty(t, testObjectStorageAPI.S3Keys)
			assert.Equal(t, "AKIAIOSFODNN7EXAMPLE", testObjectStorageAPI.S3Keys[0].AccessKey)
			assert.Equal(t, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", testObjectStorageAPI.S3Keys[0].SecretKey)
		}
	}
}

func TestGenerateKey(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		Body       string
		StatusCode int
		Error      error
	}{
		{
			Body:       `{"access_key":"AKIAIOSFODNN7NEWKEY","secret_key":"newSecretKeyExample1234567890","user_id":123,"created_at":"2023-01-01 10:00:00","updated_at":"2023-01-01 10:00:00"}`,
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
		err := testObjectStorageAPI.GenerateKey()
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
		if err == nil {
			assert.NotNil(t, testObjectStorageAPI.S3Key)
			assert.Equal(t, "AKIAIOSFODNN7NEWKEY", testObjectStorageAPI.S3Key.AccessKey)
			assert.Equal(t, "newSecretKeyExample1234567890", testObjectStorageAPI.S3Key.SecretKey)
		}
	}
}

func TestDeleteKey(t *testing.T) {
	mockHttpClient := setupMockClient(`{}`)
	testObjectStorageAPI := ObjectStorageAPI{}
	if err := testObjectStorageAPI.Init(mockHttpClient, userAuthToken, loc); err != nil {
		t.Fatalf("failed to initialize objectstorage api: %v", err)
	}
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"access_key": "AKIAIOSFODNN7EXAMPLE",
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
		accessKey := test.RequestData["access_key"].(string)
		err := testObjectStorageAPI.DeleteKey(accessKey)
		if err != nil && test.Error != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", test.Error, err)
		}
	}
}
