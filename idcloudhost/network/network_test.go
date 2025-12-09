package network

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func setupMockClient(responseBody string) *MockClient {
	return &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}
}

func TestNetworkAPI_GetNetwork(t *testing.T) {
	// Create a sample network response
	networkResponse := `{
		"id": 1,
		"name": "test-network",
		"uuid": "test-uuid",
		"user_id": 123,
		"default": true,
		"description": "Test network",
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-01T00:00:00Z"
	}`

	mockClient := setupMockClient(networkResponse)
	networkAPI := NetworkAPI{}
	if err := networkAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize network api: %v", err)
	}

	err := networkAPI.GetNetwork("test-uuid")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if networkAPI.Network == nil {
		t.Fatal("Expected Network to be populated")
	}

	if networkAPI.Network.UUID != "test-uuid" {
		t.Errorf("Expected UUID 'test-uuid', got %s", networkAPI.Network.UUID)
	}

	if networkAPI.Network.Name != "test-network" {
		t.Errorf("Expected Name 'test-network', got %s", networkAPI.Network.Name)
	}
}

func TestNetworkAPI_ListNetworks(t *testing.T) {
	// Create a sample networks response
	networksResponse := `[
		{
			"id": 1,
			"name": "test-network1",
			"uuid": "test-uuid1",
			"user_id": 123,
			"default": true,
			"description": "Test network 1",
			"created_at": "2022-01-01T00:00:00Z",
			"updated_at": "2022-01-01T00:00:00Z"
		},
		{
			"id": 2,
			"name": "test-network2",
			"uuid": "test-uuid2",
			"user_id": 123,
			"default": false,
			"description": "Test network 2",
			"created_at": "2022-01-01T00:00:00Z",
			"updated_at": "2022-01-01T00:00:00Z"
		}
	]`

	mockClient := setupMockClient(networksResponse)
	networkAPI := NetworkAPI{}
	if err := networkAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize network api: %v", err)
	}

	err := networkAPI.ListNetworks()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(networkAPI.Networks) != 2 {
		t.Fatalf("Expected 2 networks, got %d", len(networkAPI.Networks))
	}

	if networkAPI.Networks[0].UUID != "test-uuid1" {
		t.Errorf("Expected UUID 'test-uuid1', got %s", networkAPI.Networks[0].UUID)
	}

	if networkAPI.Networks[1].UUID != "test-uuid2" {
		t.Errorf("Expected UUID 'test-uuid2', got %s", networkAPI.Networks[1].UUID)
	}
}

func TestNetworkAPI_CreateDefaultNetwork(t *testing.T) {
	// Create a sample network response
	networkResponse := `{
		"id": 1,
		"name": "default-network",
		"uuid": "default-uuid",
		"user_id": 123,
		"default": true,
		"description": "Default network",
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-01T00:00:00Z"
	}`

	mockClient := setupMockClient(networkResponse)
	networkAPI := NetworkAPI{}
	if err := networkAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize network api: %v", err)
	}

	err := networkAPI.CreateDefaultNetwork()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if networkAPI.Network == nil {
		t.Fatal("Expected Network to be populated")
	}

	if networkAPI.Network.UUID != "default-uuid" {
		t.Errorf("Expected UUID 'default-uuid', got %s", networkAPI.Network.UUID)
	}

	if !networkAPI.Network.Default {
		t.Errorf("Expected Default to be true")
	}
}

func TestNetworkAPI_DeleteNetwork(t *testing.T) {
	mockClient := setupMockClient("")
	networkAPI := NetworkAPI{}
	if err := networkAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize network api: %v", err)
	}

	// Store the original DoFunc to verify it was called with the correct URL
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true
		// Check that the URL ends with the correct UUID
		if !strings.HasSuffix(req.URL.String(), "delete-uuid") {
			t.Errorf("Expected URL to end with 'delete-uuid', got %s", req.URL.String())
		}
		return originalDoFunc(req)
	}

	err := networkAPI.DeleteNetwork("delete-uuid")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Fatal("Expected DoFunc to be called")
	}
}

func TestNetworkAPI_SetAsDefault(t *testing.T) {
	// Create a sample network response
	networkResponse := `{
		"id": 1,
		"name": "new-default-network",
		"uuid": "new-default-uuid",
		"user_id": 123,
		"default": true,
		"description": "New default network",
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-01T00:00:00Z"
	}`

	mockClient := setupMockClient(networkResponse)
	networkAPI := NetworkAPI{}
	if err := networkAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize network api: %v", err)
	}

	// Store the original DoFunc to verify it was called with the correct URL
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true
		// Check that the URL contains both the UUID and "default"
		if !strings.Contains(req.URL.String(), "new-default-uuid/default") {
			t.Errorf("Expected URL to contain 'new-default-uuid/default', got %s", req.URL.String())
		}
		return originalDoFunc(req)
	}

	err := networkAPI.SetAsDefault("new-default-uuid")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Fatal("Expected DoFunc to be called")
	}

	if networkAPI.Network == nil {
		t.Fatal("Expected Network to be populated")
	}

	if networkAPI.Network.UUID != "new-default-uuid" {
		t.Errorf("Expected UUID 'new-default-uuid', got %s", networkAPI.Network.UUID)
	}

	if !networkAPI.Network.Default {
		t.Errorf("Expected Default to be true")
	}
}

func TestNetworkAPI_UpdateNetwork(t *testing.T) {
	// Create a sample network response
	networkResponse := `{
		"id": 1,
		"name": "updated-network",
		"uuid": "update-uuid",
		"user_id": 123,
		"default": false,
		"description": "Updated network",
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-02T00:00:00Z"
	}`

	mockClient := setupMockClient(networkResponse)
	networkAPI := NetworkAPI{}
	if err := networkAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize network api: %v", err)
	}

	// Store the original DoFunc to verify it was called with the correct URL and body
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true
		// Check that the URL contains the UUID
		if !strings.Contains(req.URL.String(), "update-uuid") {
			t.Errorf("Expected URL to contain 'update-uuid', got %s", req.URL.String())
		}

		// Check that the body contains the new name
		bodyBytes, _ := io.ReadAll(req.Body)
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, "updated-network") {
			t.Errorf("Expected body to contain 'updated-network', got %s", bodyString)
		}
		// Restore the body for further processing
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		return originalDoFunc(req)
	}

	err := networkAPI.UpdateNetwork("update-uuid", "updated-network")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Fatal("Expected DoFunc to be called")
	}

	if networkAPI.Network == nil {
		t.Fatal("Expected Network to be populated")
	}

	if networkAPI.Network.UUID != "update-uuid" {
		t.Errorf("Expected UUID 'update-uuid', got %s", networkAPI.Network.UUID)
	}

	if networkAPI.Network.Name != "updated-network" {
		t.Errorf("Expected Name 'updated-network', got %s", networkAPI.Network.Name)
	}
}
