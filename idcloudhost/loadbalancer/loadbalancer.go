package loadbalancer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CreateTargetRequest represents a target in the create request
type CreateTargetRequest struct {
	TargetType string `json:"target_type"`
	TargetUUID string `json:"target_uuid"`
}

// Target represents a target in the response
type Target struct {
	CreatedAt       string `json:"created_at,omitempty"`
	TargetUUID      string `json:"target_uuid"`
	TargetType      string `json:"target_type"`
	TargetIPAddress string `json:"target_ip_address,omitempty"`
}

// CreateRuleRequest represents a forwarding rule in the create request
type CreateRuleRequest struct {
	SourcePort int `json:"source_port"`
	TargetPort int `json:"target_port"`
}

// ForwardingRuleSettings represents the settings for a forwarding rule
type ForwardingRuleSettings struct {
	ConnectionLimit    int    `json:"connection_limit,omitempty"`
	SessionPersistence string `json:"session_persistence,omitempty"`
}

// ForwardingRule represents a forwarding rule in the response
type ForwardingRule struct {
	Protocol   string                  `json:"protocol,omitempty"`
	UUID       string                  `json:"uuid,omitempty"`
	CreatedAt  string                  `json:"created_at,omitempty"`
	SourcePort int                     `json:"source_port"`
	TargetPort int                     `json:"target_port"`
	Settings   *ForwardingRuleSettings `json:"settings,omitempty"`
}

// CreateLoadBalancerRequest represents the request body for creating a load balancer
type CreateLoadBalancerRequest struct {
	ReservePublicIP  bool                  `json:"reserve_public_ip,omitempty"`
	NetworkUUID      string                `json:"network_uuid"`
	Targets          []CreateTargetRequest `json:"targets,omitempty"`
	Rules            []CreateRuleRequest   `json:"rules,omitempty"`
	DisplayName      string                `json:"display_name"`
	BillingAccountID int                   `json:"billing_account_id"`
}

// LoadBalancer represents a load balancer resource
type LoadBalancer struct {
	UUID             string           `json:"uuid,omitempty"`
	DisplayName      string           `json:"display_name,omitempty"`
	UserID           int              `json:"user_id,omitempty"`
	BillingAccountID int              `json:"billing_account_id,omitempty"`
	CreatedAt        string           `json:"created_at,omitempty"`
	UpdatedAt        string           `json:"updated_at,omitempty"`
	IsDeleted        bool             `json:"is_deleted,omitempty"`
	DeletedAt        *string          `json:"deleted_at,omitempty"`
	PrivateAddress   string           `json:"private_address,omitempty"`
	NetworkUUID      string           `json:"network_uuid,omitempty"`
	ForwardingRules  []ForwardingRule `json:"forwarding_rules,omitempty"`
	Targets          []Target         `json:"targets,omitempty"`
}

type LoadBalancerAPI struct {
	c                HTTPClient
	AuthToken        string
	Location         string
	NetworkUUID      string
	ApiEndpoint      string
	LoadBalancer     *LoadBalancer
	LoadBalancerList []LoadBalancer
	Target           *Target
	ForwardingRule   *ForwardingRule
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
	req.Header.Set("apikey", lb.AuthToken)
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
	req.Header.Set("apikey", lb.AuthToken)
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
func (lb *LoadBalancerAPI) CreateLoadBalancer(request *CreateLoadBalancerRequest) error {
	payloadJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", lb.ApiEndpoint, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", lb.AuthToken)

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
		bodyBytes, _ := io.ReadAll(r.Body)
		return fmt.Errorf("status %v: %s", r.StatusCode, string(bodyBytes))
	}

	return json.NewDecoder(r.Body).Decode(&lb.LoadBalancer)
}

// RenameLoadBalancer renames a load balancer
func (lb *LoadBalancerAPI) RenameLoadBalancer(uuid string, displayName string) error {
	url := fmt.Sprintf("%s/%s", lb.ApiEndpoint, uuid)

	payload := struct {
		DisplayName string `json:"display_name"`
	}{
		DisplayName: displayName,
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
	req.Header.Set("apikey", lb.AuthToken)

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
	req.Header.Set("apikey", lb.AuthToken)
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
func (lb *LoadBalancerAPI) AddTarget(uuid string, target *CreateTargetRequest) error {
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
	req.Header.Set("apikey", lb.AuthToken)

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

	return json.NewDecoder(r.Body).Decode(&lb.Target)
}

// RemoveTarget removes a target from a load balancer
func (lb *LoadBalancerAPI) RemoveTarget(uuid string, targetUUID string) error {
	url := fmt.Sprintf("%s/%s/targets/%s", lb.ApiEndpoint, uuid, targetUUID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apikey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return nil
}

// AddRule adds a forwarding rule to a load balancer
func (lb *LoadBalancerAPI) AddRule(uuid string, rule *CreateRuleRequest) error {
	url := fmt.Sprintf("%s/%s/forwarding_rules", lb.ApiEndpoint, uuid)

	payloadJSON, err := json.Marshal(rule)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", lb.AuthToken)

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

	return json.NewDecoder(r.Body).Decode(&lb.ForwardingRule)
}

// RemoveRule removes a forwarding rule from a load balancer
func (lb *LoadBalancerAPI) RemoveRule(uuid string, ruleUUID string) error {
	url := fmt.Sprintf("%s/%s/forwarding_rules/%s", lb.ApiEndpoint, uuid, ruleUUID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apikey", lb.AuthToken)

	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusNoContent {
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
	req.Header.Set("apikey", lb.AuthToken)

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
