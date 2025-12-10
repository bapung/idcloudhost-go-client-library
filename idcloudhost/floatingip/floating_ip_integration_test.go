//go:build integration

package floatingip

import (
	"net/http"
	"os"
	"strconv"
	"testing"
)

func getEnvOrSkip(t *testing.T, key string) string {
	val := os.Getenv(key)
	if val == "" {
		t.Fatalf("Environment variable %s not set", key)
	}
	return val
}

func TestFloatingIPIntegration(t *testing.T) {
	authToken := getEnvOrSkip(t, "IDCLOUDHOST_API_KEY")
	location := getEnvOrSkip(t, "IDCLOUDHOST_LOCATION")
	billingAccountStr := getEnvOrSkip(t, "IDCLOUDHOST_BILLING_ACCOUNT")

	billingAccount, err := strconv.Atoi(billingAccountStr)
	if err != nil {
		t.Fatalf("Invalid billing account: %v", err)
	}

	client := &http.Client{}
	api := FloatingIPAPI{}
	if err := api.Init(client, authToken, location); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 1. Create floating IP
	ipName := "integration-test-ip"
	if err := api.Create(ipName, billingAccount); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	ipAddress := api.FloatingIP.Address
	t.Logf("Created floating IP: %s", ipAddress)

	// Verify the IP was created
	if api.FloatingIP.Name != ipName {
		t.Errorf("Expected name %s, got %s", ipName, api.FloatingIP.Name)
	}
	if api.FloatingIP.BillingAccount != billingAccount {
		t.Errorf("Expected billing account %d, got %d", billingAccount, api.FloatingIP.BillingAccount)
	}

	// 2. Get floating IP details
	if err := api.Get(ipAddress); err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if api.FloatingIP.Address != ipAddress {
		t.Errorf("Expected address %s, got %s", ipAddress, api.FloatingIP.Address)
	}

	// 3. Update floating IP
	updatedName := "integration-test-ip-updated"
	if err := api.Update(updatedName, billingAccount, ipAddress); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if api.FloatingIP.Name != updatedName {
		t.Errorf("Expected name %s, got %s", updatedName, api.FloatingIP.Name)
	}

	// 4. Delete floating IP (cleanup)
	if err := api.Delete(ipAddress); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	t.Logf("Deleted floating IP: %s", ipAddress)
}
