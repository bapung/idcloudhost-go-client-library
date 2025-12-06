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
		t.Skipf("Environment variable %s not set", key)
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

	// 1. Create load balancer
	lb := &LoadBalancer{
		Name:           "integration-test-lb",
		BillingAccount: billingAccount,
		ForwardRules: []ForwardingRule{
			{
				Name:         "http-rule",
				Protocol:     "http",
				FrontendPort: 80,
				BackendPort:  8080,
				HealthCheck:  true,
			},
		},
	}

	if err := api.CreateLoadBalancer(lb); err != nil {
		t.Fatalf("CreateLoadBalancer failed: %v", err)
	}
	lbUUID := api.LoadBalancer.UUID
	t.Logf("Created load balancer UUID: %s", lbUUID)

	// Verify creation
	if api.LoadBalancer.Name != "integration-test-lb" {
		t.Errorf("Expected name integration-test-lb, got %s", api.LoadBalancer.Name)
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
	if api.LoadBalancer.Name != updatedName {
		t.Errorf("Expected name %s, got %s", updatedName, api.LoadBalancer.Name)
	}

	// 5. Add a forwarding rule
	newRule := &ForwardingRule{
		Name:         "https-rule",
		Protocol:     "https",
		FrontendPort: 443,
		BackendPort:  8443,
		HealthCheck:  true,
	}
	if err := api.AddRule(lbUUID, newRule); err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}
	t.Logf("Added forwarding rule")

	// Verify rule was added
	if err := api.GetLoadBalancer(lbUUID); err != nil {
		t.Fatalf("GetLoadBalancer failed after adding rule: %v", err)
	}
	if len(api.LoadBalancer.ForwardRules) < 2 {
		t.Errorf("Expected at least 2 rules, got %d", len(api.LoadBalancer.ForwardRules))
	}

	// 6. Remove the rule we just added
	if len(api.LoadBalancer.ForwardRules) > 0 {
		// Find the rule we just added
		var ruleID int
		for _, rule := range api.LoadBalancer.ForwardRules {
			if rule.Name == "https-rule" {
				ruleID = rule.ID
				break
			}
		}
		if ruleID > 0 {
			if err := api.RemoveRule(lbUUID, ruleID); err != nil {
				t.Logf("RemoveRule failed (non-fatal): %v", err)
			} else {
				t.Logf("Removed forwarding rule ID: %d", ruleID)
			}
		}
	}

	// 7. Delete load balancer (cleanup)
	if err := api.DeleteLoadBalancer(lbUUID); err != nil {
		t.Fatalf("DeleteLoadBalancer failed: %v", err)
	}
	t.Logf("Deleted load balancer UUID: %s", lbUUID)
}
