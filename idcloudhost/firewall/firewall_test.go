//go:build !integration

package firewall

import (
	"bytes"
	"encoding/json"
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
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
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
			"display_name": "test-firewall1",
			"uuid": "test-uuid1",
			"user_id": 123,
			"billing_account_id": 1200132376,
			"description": "Test firewall 1",
			"rules": [
				{
					"direction": "inbound",
					"protocol": "tcp",
					"port_start": 22,
					"endpoint_spec_type": "any",
					"endpoint_spec": [],
					"description": "Allow SSH"
				}
			],
			"created_at": "2022-01-01T00:00:00Z",
			"updated_at": "2022-01-01T00:00:00Z"
		},
		{
			"id": 2,
			"display_name": "test-firewall2",
			"uuid": "test-uuid2",
			"user_id": 123,
			"billing_account_id": 1200132376,
			"description": "Test firewall 2",
			"rules": [
				{
					"direction": "inbound",
					"protocol": "tcp",
					"port_start": 80,
					"endpoint_spec_type": "any",
					"endpoint_spec": [],
					"description": "Allow HTTP"
				}
			],
			"created_at": "2022-01-01T00:00:00Z",
			"updated_at": "2022-01-01T00:00:00Z"
		}
	]`

	mockClient := setupMockClient(firewallsResponse)
	firewallAPI := FirewallAPI{}
	if err := firewallAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize firewall api: %v", err)
	}

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
		"display_name": "test-firewall",
		"uuid": "test-uuid",
		"user_id": 123,
		"billing_account_id": 1200132376,
		"description": "Test firewall",
		"rules": [
			{
				"direction": "inbound",
				"protocol": "tcp",
				"port_start": 22,
				"endpoint_spec_type": "any",
				"endpoint_spec": [],
				"description": "Allow SSH"
			}
		],
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-01T00:00:00Z"
	}`

	mockClient := setupMockClient(firewallResponse)
	firewallAPI := FirewallAPI{}
	if err := firewallAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize firewall api: %v", err)
	}

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
		bodyBytes, _ := io.ReadAll(req.Body)
		if err := json.Unmarshal(bodyBytes, &firewall); err != nil {
			t.Fatalf("failed to unmarshal firewall: %v", err)
		}

		if firewall.DisplayName != "test-firewall" {
			t.Errorf("Expected display_name 'test-firewall', got %s", firewall.DisplayName)
		}

		if firewall.BillingAccountID != 1200132376 {
			t.Errorf("Expected billing_account_id 1200132376, got %d", firewall.BillingAccountID)
		}

		if len(firewall.Rules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(firewall.Rules))
		}

		// Restore the body for further processing
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		return originalDoFunc(req)
	}

	// Create a firewall object to test with
	firewall := &Firewall{
		DisplayName:      "test-firewall",
		BillingAccountID: 1200132376,
		Description:      "Test firewall",
		Rules: []FirewallRule{
			{
				Direction:        "inbound",
				Protocol:         "tcp",
				PortStart:        22,
				EndpointSpecType: "any",
				EndpointSpec:     []string{},
				Description:      "Allow SSH",
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

func TestFirewallAPI_CreateFirewall_WithMultipleRules(t *testing.T) {
	// Create a sample firewall response matching the captured request
	firewallResponse := `{
		"id": 1,
		"display_name": "opopo",
		"uuid": "test-uuid",
		"user_id": 123,
		"billing_account_id": 1200132376,
		"rules": [
			{
				"direction": "inbound",
				"protocol": "any",
				"port_start": 0,
				"endpoint_spec_type": "any",
				"endpoint_spec": []
			},
			{
				"direction": "outbound",
				"protocol": "any",
				"port_start": 0,
				"endpoint_spec_type": "any",
				"endpoint_spec": []
			}
		],
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-01T00:00:00Z"
	}`

	mockClient := setupMockClient(firewallResponse)
	firewallAPI := FirewallAPI{}
	if err := firewallAPI.Init(mockClient, "test-token", "jkt03"); err != nil {
		t.Fatalf("failed to initialize firewall api: %v", err)
	}

	// Store the original DoFunc to verify it was called with the correct body
	originalDoFunc := mockClient.DoFunc
	called := false

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		called = true

		// Decode the request body to verify it matches the captured request structure
		var firewall Firewall
		bodyBytes, _ := io.ReadAll(req.Body)
		if err := json.Unmarshal(bodyBytes, &firewall); err != nil {
			t.Fatalf("failed to unmarshal firewall: %v", err)
		}

		if firewall.DisplayName != "opopo" {
			t.Errorf("Expected display_name 'opopo', got %s", firewall.DisplayName)
		}

		if firewall.BillingAccountID != 1200132376 {
			t.Errorf("Expected billing_account_id 1200132376, got %d", firewall.BillingAccountID)
		}

		if len(firewall.Rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(firewall.Rules))
		}

		// Verify first rule (inbound)
		if firewall.Rules[0].Direction != "inbound" {
			t.Errorf("Expected first rule direction 'inbound', got %s", firewall.Rules[0].Direction)
		}
		if firewall.Rules[0].Protocol != "any" {
			t.Errorf("Expected first rule protocol 'any', got %s", firewall.Rules[0].Protocol)
		}
		if firewall.Rules[0].EndpointSpecType != "any" {
			t.Errorf("Expected first rule endpoint_spec_type 'any', got %s", firewall.Rules[0].EndpointSpecType)
		}

		// Verify second rule (outbound)
		if firewall.Rules[1].Direction != "outbound" {
			t.Errorf("Expected second rule direction 'outbound', got %s", firewall.Rules[1].Direction)
		}
		if firewall.Rules[1].Protocol != "any" {
			t.Errorf("Expected second rule protocol 'any', got %s", firewall.Rules[1].Protocol)
		}

		// Restore the body for further processing
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		return originalDoFunc(req)
	}

	// Create a firewall object matching the captured request
	firewall := &Firewall{
		DisplayName:      "opopo",
		BillingAccountID: 1200132376,
		Rules: []FirewallRule{
			{
				Direction:        "inbound",
				Protocol:         "any",
				PortStart:        0,
				EndpointSpecType: "any",
				EndpointSpec:     []string{},
			},
			{
				Direction:        "outbound",
				Protocol:         "any",
				PortStart:        0,
				EndpointSpecType: "any",
				EndpointSpec:     []string{},
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
}

func TestFirewallAPI_UpdateFirewall(t *testing.T) {
	// Create a sample firewall response
	firewallResponse := `{
		"id": 1,
		"display_name": "updated-firewall",
		"uuid": "update-uuid",
		"user_id": 123,
		"billing_account_id": 1200132376,
		"description": "Updated firewall",
		"rules": [
			{
				"direction": "inbound",
				"protocol": "tcp",
				"port_start": 80,
				"endpoint_spec_type": "any",
				"endpoint_spec": [],
				"description": "Allow HTTP"
			}
		],
		"created_at": "2022-01-01T00:00:00Z",
		"updated_at": "2022-01-02T00:00:00Z"
	}`

	mockClient := setupMockClient(firewallResponse)
	firewallAPI := FirewallAPI{}
	if err := firewallAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize firewall api: %v", err)
	}

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

		// Decode the request body to verify it - UpdateFirewall only sends rules
		var payload struct {
			Rules []FirewallRule `json:"rules"`
		}
		bodyBytes, _ := io.ReadAll(req.Body)
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			t.Fatalf("failed to unmarshal payload: %v", err)
		}

		if len(payload.Rules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(payload.Rules))
		}

		// Restore the body for further processing
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		return originalDoFunc(req)
	}

	// Create a firewall object to test with
	firewall := &Firewall{
		DisplayName:      "updated-firewall",
		BillingAccountID: 1200132376,
		Description:      "Updated firewall",
		Rules: []FirewallRule{
			{
				Direction:        "inbound",
				Protocol:         "tcp",
				PortStart:        80,
				EndpointSpecType: "any",
				EndpointSpec:     []string{},
				Description:      "Allow HTTP",
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

	if firewallAPI.Firewall.DisplayName != "updated-firewall" {
		t.Errorf("Expected DisplayName 'updated-firewall', got %s", firewallAPI.Firewall.DisplayName)
	}
}

func TestFirewallAPI_DeleteFirewall(t *testing.T) {
	mockClient := setupMockClient("{}")
	firewallAPI := FirewallAPI{}
	if err := firewallAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize firewall api: %v", err)
	}

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
	if err := firewallAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize firewall api: %v", err)
	}

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
		bodyBytes, _ := io.ReadAll(req.Body)
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, "vm-uuid") {
			t.Errorf("Expected body to contain 'vm-uuid', got %s", bodyString)
		}

		// Restore the body for further processing
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

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
	if err := firewallAPI.Init(mockClient, "test-token", "test-location"); err != nil {
		t.Fatalf("failed to initialize firewall api: %v", err)
	}

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
		bodyBytes, _ := io.ReadAll(req.Body)
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, "vm-uuid") {
			t.Errorf("Expected body to contain 'vm-uuid', got %s", bodyString)
		}

		// Restore the body for further processing
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

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

func TestValidateFirewallRules(t *testing.T) {
	tests := []struct {
		name    string
		rules   []FirewallRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid inbound tcp rule",
			rules: []FirewallRule{
				{
					Direction:        "inbound",
					Protocol:         "tcp",
					PortStart:        22,
					EndpointSpecType: "any",
					EndpointSpec:     []string{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid outbound any rule",
			rules: []FirewallRule{
				{
					Direction:        "outbound",
					Protocol:         "any",
					PortStart:        0,
					EndpointSpecType: "any",
					EndpointSpec:     []string{},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid direction",
			rules: []FirewallRule{
				{
					Direction:        "invalid",
					Protocol:         "tcp",
					PortStart:        22,
					EndpointSpecType: "any",
					EndpointSpec:     []string{},
				},
			},
			wantErr: true,
			errMsg:  "direction must be either 'inbound' or 'outbound'",
		},
		{
			name: "invalid protocol",
			rules: []FirewallRule{
				{
					Direction:        "inbound",
					Protocol:         "invalid",
					PortStart:        22,
					EndpointSpecType: "any",
					EndpointSpec:     []string{},
				},
			},
			wantErr: true,
			errMsg:  "protocol must be one of 'tcp', 'udp', 'icmp', or 'any'",
		},
		{
			name: "invalid endpoint_spec_type",
			rules: []FirewallRule{
				{
					Direction:        "inbound",
					Protocol:         "tcp",
					PortStart:        22,
					EndpointSpecType: "invalid",
					EndpointSpec:     []string{},
				},
			},
			wantErr: true,
			errMsg:  "endpoint_spec_type must be one of 'any', 'cidr', 'firewall', or 'ip_prefixes'",
		},
		{
			name: "nil endpoint_spec",
			rules: []FirewallRule{
				{
					Direction:        "inbound",
					Protocol:         "tcp",
					PortStart:        22,
					EndpointSpecType: "any",
					EndpointSpec:     nil,
				},
			},
			wantErr: true,
			errMsg:  "endpoint_spec cannot be nil",
		},
		{
			name: "valid ip_prefixes endpoint_spec_type",
			rules: []FirewallRule{
				{
					Direction:        "inbound",
					Protocol:         "tcp",
					PortStart:        22,
					EndpointSpecType: "ip_prefixes",
					EndpointSpec:     []string{"10.90.0.0/24"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid rule with port range",
			rules: []FirewallRule{
				{
					Direction:        "inbound",
					Protocol:         "tcp",
					PortStart:        8000,
					PortEnd:          8088,
					EndpointSpecType: "any",
					EndpointSpec:     []string{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid rule with uuid (for updates)",
			rules: []FirewallRule{
				{
					UUID:             "7ef7cf38-5a3a-4a77-ba48-85750529df62",
					Direction:        "inbound",
					Protocol:         "tcp",
					PortStart:        22,
					EndpointSpecType: "any",
					EndpointSpec:     []string{},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFirewallRules(tt.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFirewallRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateFirewallRules() error = %v, expected error containing %v", err, tt.errMsg)
			}
		})
	}
}
