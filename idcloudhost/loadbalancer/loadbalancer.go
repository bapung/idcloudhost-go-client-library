package loadbalancer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Target struct {
	VMID   int    `json:"vm_id,omitempty"`
	VMUUID string `json:"vm_uuid,omitempty"`
	IPAddr string `json:"ip_addr,omitempty"`
	Name   string `json:"name,omitempty"`
}

type ForwardingRule struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	FrontendPort  int    `json:"frontend_port"`
	BackendPort   int    `json:"backend_port"`
	HealthCheck   bool   `json:"health_check,omitempty"`
	HealthTimeout int    `json:"health_timeout,omitempty"`
}

type LoadBalancer struct {
	ID             int              `json:"id,omitempty"`
	UUID           string           `json:"uuid,omitempty"`
	Name           string           `json:"name"`
	BillingAccount int              `json:"billing_account_id"`
	UserID         int              `json:"user_id,omitempty"`
	TargetIPs      []Target         `json:"target_ips,omitempty"`
	ForwardRules   []ForwardingRule `json:"forward_rules,omitempty"`
	CreatedAt      string           `json:"created_at,omitempty"`
	UpdatedAt      string           `json:"updated_at,omitempty"`
}

type LoadBalancerAPI struct {
	c                HTTPClient
	AuthToken        string
	Location         string
	ApiEndpoint      string
	LoadBalancer     *LoadBalancer
	LoadBalancerList []LoadBalancer
}

func (lb *LoadBalancerAPI) Init(c HTTPClient, authToken string, location string) error {
	lb.c = c
	lb.AuthToken = authToken
	lb.Location = location
	lb.ApiEndpoint = fmt.Sprintf(
		"https://api.idcloudhost.com/v1/%s/network/load_balancers",
		lb.Location,
	)
	r, err := http.Get(lb.ApiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("location: %s not found", lb.Location)
	}
	return nil
}

// ListLoadBalancers lists all load balancers for the user
func (lb *LoadBalancerAPI) ListLoadBalancers() error {
	req, err := http.NewRequest("GET", lb.ApiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", lb.AuthToken)
	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	return json.NewDecoder(r.Body).Decode(&lb.LoadBalancerList)
}

// GetLoadBalancer retrieves information about a specific load balancer
func (lb *LoadBalancerAPI) GetLoadBalancer(uuid string) error {
	url := fmt.Sprintf("%s/%s", lb.ApiEndpoint, uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", lb.AuthToken)
	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	return json.NewDecoder(r.Body).Decode(&lb.LoadBalancer)
}

// CreateLoadBalancer creates a new load balancer
func (lb *LoadBalancerAPI) CreateLoadBalancer(loadBalancer *LoadBalancer) error {
	payloadJSON, err := json.Marshal(loadBalancer)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", lb.ApiEndpoint, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&lb.LoadBalancer)
}

// RenameLoadBalancer renames a load balancer
func (lb *LoadBalancerAPI) RenameLoadBalancer(uuid string, name string) error {
	url := fmt.Sprintf("%s/%s", lb.ApiEndpoint, uuid)

	payload := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&lb.LoadBalancer)
}

// DeleteLoadBalancer deletes a load balancer
func (lb *LoadBalancerAPI) DeleteLoadBalancer(uuid string) error {
	url := fmt.Sprintf("%s/%s", lb.ApiEndpoint, uuid)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", lb.AuthToken)
	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	return nil
}

// AddTarget adds a target to a load balancer
func (lb *LoadBalancerAPI) AddTarget(uuid string, target *Target) error {
	url := fmt.Sprintf("%s/%s/targets", lb.ApiEndpoint, uuid)

	payloadJSON, err := json.Marshal(target)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&lb.LoadBalancer)
}

// RemoveTarget removes a target from a load balancer
func (lb *LoadBalancerAPI) RemoveTarget(uuid string, vmUUID string) error {
	url := fmt.Sprintf("%s/%s/targets", lb.ApiEndpoint, uuid)

	payload := struct {
		VMUUID string `json:"vm_uuid"`
	}{
		VMUUID: vmUUID,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return nil
}

// AddRule adds a forwarding rule to a load balancer
func (lb *LoadBalancerAPI) AddRule(uuid string, rule *ForwardingRule) error {
	url := fmt.Sprintf("%s/%s/rules", lb.ApiEndpoint, uuid)

	payloadJSON, err := json.Marshal(rule)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&lb.LoadBalancer)
}

// RemoveRule removes a forwarding rule from a load balancer
func (lb *LoadBalancerAPI) RemoveRule(uuid string, ruleID int) error {
	url := fmt.Sprintf("%s/%s/rules/%d", lb.ApiEndpoint, uuid, ruleID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return nil
}

// ChangeBillingAccount changes the billing account for a load balancer
func (lb *LoadBalancerAPI) ChangeBillingAccount(uuid string, billingAccountID int) error {
	url := fmt.Sprintf("%s/%s/billing_account", lb.ApiEndpoint, uuid)

	payload := struct {
		BillingAccountID int `json:"billing_account_id"`
	}{
		BillingAccountID: billingAccountID,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&lb.LoadBalancer)
}
