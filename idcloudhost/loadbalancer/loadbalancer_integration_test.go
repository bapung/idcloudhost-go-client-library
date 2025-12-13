//go:build integration

package loadbalancer

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

func TestLoadBalancerIntegration(t *testing.T) {
	authToken := getEnvOrSkip(t, "IDCLOUDHOST_API_KEY")
	location := getEnvOrSkip(t, "IDCLOUDHOST_LOCATION")
	billingAccountStr := getEnvOrSkip(t, "IDCLOUDHOST_BILLING_ACCOUNT")

	billingAccount, err := strconv.Atoi(billingAccountStr)
	if err != nil {
		t.Fatalf("Invalid billing account: %v", err)
	}

	client := &http.Client{}
	api := LoadBalancerAPI{}
	if err := api.Init(client, authToken, location); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Set the network UUID for load balancer creation
	api.NetworkUUID = "730bf645-9b36-44b8-8ca1-46d2480cc0d6"

	// 1. Create load balancer
	targetVMUUID := "7171a7ea-507e-4e78-9d94-6c15abeb5f32"
	createReq := &CreateLoadBalancerRequest{
		ReservePublicIP:  true,
		NetworkUUID:      api.NetworkUUID,
		DisplayName:      "integration-test-lb",
		BillingAccountID: billingAccount,
		Targets: []CreateTargetRequest{
			{
				TargetType: "vm",
				TargetUUID: targetVMUUID,
			},
		},
		Rules: []CreateRuleRequest{
			{
				SourcePort: 80,
				TargetPort: 8080,
			},
		},
	}

	t.Logf("Creating load balancer with NetworkUUID: %s, BillingAccountID: %d", api.NetworkUUID, billingAccount)
	if err := api.CreateLoadBalancer(createReq); err != nil {
		t.Fatalf("CreateLoadBalancer failed: %v", err)
	}
	lbUUID := api.LoadBalancer.UUID
	t.Logf("Created load balancer UUID: %s", lbUUID)
	t.Logf("Private address: %s", api.LoadBalancer.PrivateAddress)

	// Verify creation
	if api.LoadBalancer.DisplayName != "integration-test-lb" {
		t.Errorf("Expected display_name integration-test-lb, got %s", api.LoadBalancer.DisplayName)
	}
	if api.LoadBalancer.NetworkUUID != api.NetworkUUID {
		t.Errorf("Expected network_uuid %s, got %s", api.NetworkUUID, api.LoadBalancer.NetworkUUID)
	}

	// 2. List load balancers and verify existence
	if err := api.ListLoadBalancers(); err != nil {
		t.Fatalf("ListLoadBalancers failed: %v", err)
	}
	found := false
	for _, loadBalancer := range api.LoadBalancerList {
		if loadBalancer.UUID == lbUUID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Created load balancer not found in list")
	}

	// 3. Get load balancer details
	if err := api.GetLoadBalancer(lbUUID); err != nil {
		t.Fatalf("GetLoadBalancer failed: %v", err)
	}
	if api.LoadBalancer.UUID != lbUUID {
		t.Errorf("Expected UUID %s, got %s", lbUUID, api.LoadBalancer.UUID)
	}

	// 4. Rename load balancer
	updatedName := "integration-test-lb-updated"
	if err := api.RenameLoadBalancer(lbUUID, updatedName); err != nil {
		t.Fatalf("RenameLoadBalancer failed: %v", err)
	}
	if api.LoadBalancer.DisplayName != updatedName {
		t.Errorf("Expected display_name %s, got %s", updatedName, api.LoadBalancer.DisplayName)
	}

	// 5. Add a forwarding rule
	newRule := &CreateRuleRequest{
		SourcePort: 443,
		TargetPort: 8443,
	}
	if err := api.AddRule(lbUUID, newRule); err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	t.Logf("Added forwarding rule")

	// Verify rule was added
	if api.ForwardingRule != nil {
		t.Logf("Rule UUID: %s, Protocol: %s", api.ForwardingRule.UUID, api.ForwardingRule.Protocol)
	}
	if err := api.GetLoadBalancer(lbUUID); err != nil {
		t.Fatalf("GetLoadBalancer failed after adding rule: %v", err)
	}
	if len(api.LoadBalancer.ForwardingRules) < 2 {
		t.Errorf("Expected at least 2 rules, got %d", len(api.LoadBalancer.ForwardingRules))
	}

	// 6. Remove the rule we just added
	if len(api.LoadBalancer.ForwardingRules) > 0 {
		// Find the rule we just added (443 -> 8443)
		var ruleUUID string
		for _, rule := range api.LoadBalancer.ForwardingRules {
			if rule.SourcePort == 443 && rule.TargetPort == 8443 {
				ruleUUID = rule.UUID
				break
			}
		}
		if ruleUUID != "" {
			if err := api.RemoveRule(lbUUID, ruleUUID); err != nil {
				t.Logf("RemoveRule failed (non-fatal): %v", err)
			} else {
				t.Logf("Removed forwarding rule UUID: %s", ruleUUID)
			}
		}
	}

	// Test adding an additional target (we already have one from create)
	anotherVMUUID := "70517643-b046-48e3-9bae-83dc2c143beb"
	t.Logf("Testing target operations with VM UUID: %s", anotherVMUUID)

	// Add another target
	target := &CreateTargetRequest{
		TargetType: "vm",
		TargetUUID: anotherVMUUID,
	}
	if err := api.AddTarget(lbUUID, target); err != nil {
		t.Logf("AddTarget failed (non-fatal): %v", err)
	} else {
		t.Logf("Added target VM: %s", anotherVMUUID)
		if api.Target != nil {
			t.Logf("Target created_at: %s, target_ip_address: %s", api.Target.CreatedAt, api.Target.TargetIPAddress)
		}

		// Remove the additional target
		if err := api.RemoveTarget(lbUUID, anotherVMUUID); err != nil {
			t.Logf("RemoveTarget failed (non-fatal): %v", err)
		} else {
			t.Logf("Removed target VM: %s", anotherVMUUID)
		}
	}

	// 7. Delete load balancer (cleanup)
	if err := api.DeleteLoadBalancer(lbUUID); err != nil {
		t.Fatalf("DeleteLoadBalancer failed: %v", err)
	}
	t.Logf("Deleted load balancer UUID: %s", lbUUID)

	// Note: If a public IP was reserved (reserve_public_ip: true), you may need to manually
	// delete the floating IP from the console after the test completes.
	// The floating IP is not automatically deleted when the load balancer is deleted.
	// To clean up manually:
	// 1. Go to IDCloudHost Console -> Network -> Floating IPs
	// 2. Find unassigned IPs created around the test time
	// 3. Delete them manually or use the cleanup script:
	//    go run scripts/cleanup_loadbalancer_ips.go -apikey=$IDCLOUDHOST_API_KEY -location=$IDCLOUDHOST_LOCATION delete <IP_ADDRESS>
}
