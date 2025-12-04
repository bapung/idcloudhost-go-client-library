package firewall

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
				Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}
}

func TestFirewallAPI_ListFirewalls(t *testing.T) {
	// Create a sample firewalls response
	firewallsResponse := `[
		{
			"id": 1,
			"name": "test-firewall1",
			"uuid": "test-uuid1",
			"user_id": 123,
			"description": "Test firewall 1",
			"rules": [
				{
					"type": "ingress",
					"protocol": "tcp",
					"port_range": "22",
					"source": "0.0.0.0/0",
					"description": "Allow SSH"
				}
			],
			"created_at": "2022-01-01T00:00:00Z",
			"updated_at": "2022-01-01T00:00:00Z"
		},
		{
			"id": 2,
			"name": "test-firewall2",
			"uuid": "test-uuid2",
			"user_id": 123,
			"description": "Test firewall 2",
			"rules": [
				{
					"type": "ingress",
					"protocol": "tcp",
					"port_range": "80",
					"source": "0.0.0.0/0",
					"description": "Allow HTTP"
				}
			],
			"created_at": "2022-01-01T00:00:00Z",
			"updated_at": "2022-01-01T00:00:00Z"
		}
	]`

	mockClient := setupMockClient(firewallsResponse)
	firewallAPI := FirewallAPI{}
	firewallAPI.Init(mockClient, "test-token", "test-location")

	err := firewallAPI.ListFirewalls()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(firewallAPI.Firewalls) != 2 {
		t.Fatalf("Expected 2 firewalls, got %d", len(firewallAPI.Firewalls))
	}

	if firewallAPI.Firewalls[0].UUID != "test-uuid1" {
		t.Errorf("Expected UUID 'test-uuid1', got %s", firewallAPI.Firewalls[0].UUID)
	}

	if firewallAPI.Firewalls[1].UUID != "test-uuid2" {
		t.Errorf("Expected UUID 'test-uuid2', got %s", firewallAPI.Firewalls[1].UUID)
	}
}

func TestFirewallAPI_CreateFirewall(t *testing.T) {
	// Create a sample firewall response
	firewallResponse := `{
		"id": 1,
		"name": "test-firewall",
		"uuid": "test-uuid",
		"user_id": 123,
		"description": "Test firewall",
		"rules": [
			{
				"type": "ingress",
				"protocol": "tcp",
				"port_range": "22",
				"source": "0.0.0.0/0",
				"description": "Allow SSH"
			}
		],
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-01T00:00:00Z"
	}`

	mockClient := setupMockClient(firewallResponse)
	firewallAPI := FirewallAPI{}
	firewallAPI.Init(mockClient, "test-token", "test-location")

	// Store the original DoFunc to verify it was called with the correct body
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true

		// Check the method
		if req.Method != "POST" {
			t.Errorf("Expected method POST, got %s", req.Method)
		}

		// Check the content type
		contentType := req.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
		}

		// Decode the request body to verify it
		var firewall Firewall
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(bodyBytes, &firewall)

		if firewall.Name != "test-firewall" {
			t.Errorf("Expected name 'test-firewall', got %s", firewall.Name)
		}

		if len(firewall.Rules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(firewall.Rules))
		}

		// Restore the body for further processing
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		return originalDoFunc(req)
	}

	// Create a firewall object to test with
	firewall := &Firewall{
		Name:        "test-firewall",
		Description: "Test firewall",
		Rules: []FirewallRule{
			{
				Type:        "ingress",
				Protocol:    "tcp",
				PortRange:   "22",
				Source:      "0.0.0.0/0",
				Description: "Allow SSH",
			},
		},
	}

	err := firewallAPI.CreateFirewall(firewall)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Fatal("Expected DoFunc to be called")
	}

	if firewallAPI.Firewall == nil {
		t.Fatal("Expected Firewall to be populated")
	}

	if firewallAPI.Firewall.UUID != "test-uuid" {
		t.Errorf("Expected UUID 'test-uuid', got %s", firewallAPI.Firewall.UUID)
	}

	if len(firewallAPI.Firewall.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(firewallAPI.Firewall.Rules))
	}
}

func TestFirewallAPI_UpdateFirewall(t *testing.T) {
	// Create a sample firewall response
	firewallResponse := `{
		"id": 1,
		"name": "updated-firewall",
		"uuid": "update-uuid",
		"user_id": 123,
		"description": "Updated firewall",
		"rules": [
			{
				"type": "ingress",
				"protocol": "tcp",
				"port_range": "80",
				"source": "0.0.0.0/0",
				"description": "Allow HTTP"
			}
		],
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-02T00:00:00Z"
	}`

	mockClient := setupMockClient(firewallResponse)
	firewallAPI := FirewallAPI{}
	firewallAPI.Init(mockClient, "test-token", "test-location")

	// Store the original DoFunc to verify it was called with the correct URL and body
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true

		// Check the method
		if req.Method != "PUT" {
			t.Errorf("Expected method PUT, got %s", req.Method)
		}

		// Check that the URL contains the UUID
		if !strings.Contains(req.URL.String(), "update-uuid") {
			t.Errorf("Expected URL to contain 'update-uuid', got %s", req.URL.String())
		}

		// Check the content type
		contentType := req.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
		}

		// Decode the request body to verify it
		var firewall Firewall
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(bodyBytes, &firewall)

		if firewall.Name != "updated-firewall" {
			t.Errorf("Expected name 'updated-firewall', got %s", firewall.Name)
		}

		if len(firewall.Rules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(firewall.Rules))
		}

		// Restore the body for further processing
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		return originalDoFunc(req)
	}

	// Create a firewall object to test with
	firewall := &Firewall{
		Name:        "updated-firewall",
		Description: "Updated firewall",
		Rules: []FirewallRule{
			{
				Type:        "ingress",
				Protocol:    "tcp",
				PortRange:   "80",
				Source:      "0.0.0.0/0",
				Description: "Allow HTTP",
			},
		},
	}

	err := firewallAPI.UpdateFirewall("update-uuid", firewall)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Fatal("Expected DoFunc to be called")
	}

	if firewallAPI.Firewall == nil {
		t.Fatal("Expected Firewall to be populated")
	}

	if firewallAPI.Firewall.UUID != "update-uuid" {
		t.Errorf("Expected UUID 'update-uuid', got %s", firewallAPI.Firewall.UUID)
	}

	if firewallAPI.Firewall.Name != "updated-firewall" {
		t.Errorf("Expected Name 'updated-firewall', got %s", firewallAPI.Firewall.Name)
	}
}

func TestFirewallAPI_DeleteFirewall(t *testing.T) {
	mockClient := setupMockClient("{}")
	firewallAPI := FirewallAPI{}
	firewallAPI.Init(mockClient, "test-token", "test-location")

	// Store the original DoFunc to verify it was called with the correct URL
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true

		// Check the method
		if req.Method != "DELETE" {
			t.Errorf("Expected method DELETE, got %s", req.Method)
		}

		// Check that the URL ends with the correct UUID
		if !strings.HasSuffix(req.URL.String(), "delete-uuid") {
			t.Errorf("Expected URL to end with 'delete-uuid', got %s", req.URL.String())
		}

		return originalDoFunc(req)
	}

	err := firewallAPI.DeleteFirewall("delete-uuid")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Fatal("Expected DoFunc to be called")
	}
}

func TestFirewallAPI_AssignFirewall(t *testing.T) {
	mockClient := setupMockClient("{}")
	firewallAPI := FirewallAPI{}
	firewallAPI.Init(mockClient, "test-token", "test-location")

	// Store the original DoFunc to verify it was called with the correct URL and body
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true

		// Check the method
		if req.Method != "POST" {
			t.Errorf("Expected method POST, got %s", req.Method)
		}

		// Check that the URL contains the firewall UUID and "assign"
		if !strings.Contains(req.URL.String(), "firewall-uuid/assign") {
			t.Errorf("Expected URL to contain 'firewall-uuid/assign', got %s", req.URL.String())
		}

		// Check the content type
		contentType := req.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
		}

		// Check the body contains the VM UUID
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, "vm-uuid") {
			t.Errorf("Expected body to contain 'vm-uuid', got %s", bodyString)
		}

		// Restore the body for further processing
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		return originalDoFunc(req)
	}

	err := firewallAPI.AssignFirewall("firewall-uuid", "vm-uuid")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Fatal("Expected DoFunc to be called")
	}
}

func TestFirewallAPI_UnassignFirewall(t *testing.T) {
	mockClient := setupMockClient("{}")
	firewallAPI := FirewallAPI{}
	firewallAPI.Init(mockClient, "test-token", "test-location")

	// Store the original DoFunc to verify it was called with the correct URL and body
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true

		// Check the method
		if req.Method != "POST" {
			t.Errorf("Expected method POST, got %s", req.Method)
		}

		// Check that the URL contains the firewall UUID and "unassign"
		if !strings.Contains(req.URL.String(), "firewall-uuid/unassign") {
			t.Errorf("Expected URL to contain 'firewall-uuid/unassign', got %s", req.URL.String())
		}

		// Check the content type
		contentType := req.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
		}

		// Check the body contains the VM UUID
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, "vm-uuid") {
			t.Errorf("Expected body to contain 'vm-uuid', got %s", bodyString)
		}

		// Restore the body for further processing
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		return originalDoFunc(req)
	}

	err := firewallAPI.UnassignFirewall("firewall-uuid", "vm-uuid")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !called {
		t.Fatal("Expected DoFunc to be called")
	}
}
