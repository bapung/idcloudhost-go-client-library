//go:build integration

package firewall

import (
	"net/http"
	"os"
	"testing"
)

func getEnvOrSkip(t *testing.T, key string) string {
	val := os.Getenv(key)
	if val == "" {
		t.Skipf("Environment variable %s not set", key)
	}
	return val
}

func TestFirewallIntegration(t *testing.T) {
	authToken := getEnvOrSkip(t, "IDCLOUDHOST_API_KEY")
	location := getEnvOrSkip(t, "IDCLOUDHOST_LOCATION")

	client := &http.Client{}
	api := FirewallAPI{}
	if err := api.Init(client, authToken, location); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 1. Create firewall
	firewall := &Firewall{
		Name:        "integration-test-fw",
		Description: "Integration test firewall",
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

	// 3. Update firewall
	firewall.Name = "integration-test-fw-updated"
	if err := api.UpdateFirewall(uuid, firewall); err != nil {
		t.Fatalf("UpdateFirewall failed: %v", err)
	}
	if api.Firewall.Name != "integration-test-fw-updated" {
		t.Errorf("Firewall name not updated")
	}

	// 4. Delete firewall (cleanup)
	if err := api.DeleteFirewall(uuid); err != nil {
		t.Fatalf("DeleteFirewall failed: %v", err)
	}
	t.Logf("Deleted firewall UUID: %s", uuid)
}
