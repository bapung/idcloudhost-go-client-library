//go:build integration

package firewall

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func getEnvOrSkip(t *testing.T, key string) string {
	val := os.Getenv(key)
	if val == "" {
		t.Fatalf("Environment variable %s not set", key)
	}
	return val
}

func getEnvOrDefault(t *testing.T, key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func TestFirewallIntegration(t *testing.T) {
	authToken := getEnvOrSkip(t, "IDCLOUDHOST_API_KEY")
	location := getEnvOrSkip(t, "IDCLOUDHOST_LOCATION")
	billingAccountID := getEnvOrDefault(t, "IDCLOUDHOST_BILLING_ACCOUNT_ID", "1200132376")

	client := &http.Client{}
	api := FirewallAPI{}
	if err := api.Init(client, authToken, location); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Convert billing account ID from string to int
	var billingID int
	if _, err := fmt.Sscanf(billingAccountID, "%d", &billingID); err != nil {
		t.Fatalf("Invalid billing account ID: %v", err)
	}

	// 1. Create firewall with new API structure
	firewall := &Firewall{
		DisplayName:      "integration-test-fw",
		BillingAccountID: billingID,
		Description:      "Integration test firewall",
		Rules: []FirewallRule{
			{
				Direction:        "inbound",
				Protocol:         "tcp",
				PortStart:        22,
				EndpointSpecType: "any",
				EndpointSpec:     []string{},
				Description:      "Allow SSH",
			},
			{
				Direction:        "outbound",
				Protocol:         "any",
				PortStart:        0,
				EndpointSpecType: "any",
				EndpointSpec:     []string{},
				Description:      "Allow all outbound",
			},
		},
	}
	if err := api.CreateFirewall(firewall); err != nil {
		t.Fatalf("CreateFirewall failed: %v", err)
	}
	uuid := api.Firewall.UUID
	t.Logf("Created firewall UUID: %s", uuid)

	// 2. List firewalls and check existence
	if err := api.ListFirewalls(); err != nil {
		t.Fatalf("ListFirewalls failed: %v", err)
	}
	found := false
	for _, fw := range api.Firewalls {
		if fw.UUID == uuid {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Created firewall not found in list")
	}

	// 3. Update firewall - preserve rule UUIDs from created firewall
	// Get the rules with UUIDs from the created firewall
	updatedFirewall := &Firewall{
		Rules: api.Firewall.Rules, // Preserve existing rules with their UUIDs
	}
	// Add a new rule without UUID
	updatedFirewall.Rules = append(updatedFirewall.Rules, FirewallRule{
		Direction:        "inbound",
		Protocol:         "tcp",
		PortStart:        80,
		EndpointSpecType: "any",
		EndpointSpec:     []string{},
		Description:      "Allow HTTP",
	})

	if err := api.UpdateFirewall(uuid, updatedFirewall); err != nil {
		t.Fatalf("UpdateFirewall failed: %v", err)
	}

	// Verify the update
	if len(api.Firewall.Rules) != 3 {
		t.Errorf("Expected 3 rules after update, got %d", len(api.Firewall.Rules))
	}

	// 4. Delete firewall (cleanup)
	if err := api.DeleteFirewall(uuid); err != nil {
		t.Fatalf("DeleteFirewall failed: %v", err)
	}
	t.Logf("Deleted firewall UUID: %s", uuid)
}
