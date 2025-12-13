//go:build integration

package network

import (
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

func TestNetworkIntegration(t *testing.T) {
	authToken := getEnvOrSkip(t, "IDCLOUDHOST_API_KEY")
	location := getEnvOrSkip(t, "IDCLOUDHOST_LOCATION")

	client := &http.Client{}
	api := NetworkAPI{}
	if err := api.Init(client, authToken, location); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 1. Create network
	if err := api.CreateNetwork("integration-test-network"); err != nil {
		t.Fatalf("CreateNetwork failed: %v", err)
	}
	networkUUID := api.Network.UUID
	t.Logf("Created network UUID: %s", networkUUID)

	// Verify the network was created
	if api.Network.UUID == "" {
		t.Error("Expected non-empty UUID")
	}

	// 2. List networks and verify existence
	if err := api.ListNetworks(); err != nil {
		t.Fatalf("ListNetworks failed: %v", err)
	}
	found := false
	for _, net := range api.Networks {
		if net.UUID == networkUUID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Created network not found in list")
	}

	// 3. Get network details
	if err := api.GetNetwork(networkUUID); err != nil {
		t.Fatalf("GetNetwork failed: %v", err)
	}
	if api.Network.UUID != networkUUID {
		t.Errorf("Expected UUID %s, got %s", networkUUID, api.Network.UUID)
	}

	// 4. Update network name
	updatedName := "integration-test-network-updated"
	if err := api.UpdateNetwork(networkUUID, updatedName); err != nil {
		t.Fatalf("UpdateNetwork failed: %v", err)
	}
	if api.Network.Name != updatedName {
		t.Errorf("Expected name %s, got %s", updatedName, api.Network.Name)
	}
	// default cannot be deleted
	// 5. Set as default (if not already)
	//if err := api.SetAsDefault(networkUUID); err != nil {
	//	t.Fatalf("SetAsDefault failed: %v", err)
	//}

	// 6. Delete network (cleanup)
	if err := api.DeleteNetwork(networkUUID); err != nil {
		t.Fatalf("DeleteNetwork failed: %v", err)
	}
	t.Logf("Deleted network UUID: %s", networkUUID)
}
