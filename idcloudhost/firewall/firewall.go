package firewall

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

type FirewallRule struct {
	UUID             string                 `json:"uuid,omitempty"`        // Rule UUID (for updates)
	Type             map[string]interface{} `json:"type,omitempty"`        // UI metadata (optional, can be ignored)
	Direction        string                 `json:"direction"`             // "inbound" or "outbound"
	Protocol         string                 `json:"protocol"`              // "tcp", "udp", "icmp", "any"
	PortStart        int                    `json:"port_start,omitempty"`  // Starting port number
	PortEnd          int                    `json:"port_end,omitempty"`    // Ending port number (optional)
	EndpointSpecType string                 `json:"endpoint_spec_type"`    // "any", "cidr", "firewall", "ip_prefixes"
	EndpointSpec     []string               `json:"endpoint_spec"`         // List of CIDR blocks or firewall IDs
	Description      string                 `json:"description,omitempty"` // Optional description
}

type Firewall struct {
	ID               int            `json:"id,omitempty"`
	DisplayName      string         `json:"display_name"`       // Use display_name as per API
	BillingAccountID int            `json:"billing_account_id"` // Required for creation
	UUID             string         `json:"uuid,omitempty"`
	UserID           int            `json:"user_id,omitempty"`
	Description      string         `json:"description,omitempty"`
	Rules            []FirewallRule `json:"rules"`
	CreatedAt        string         `json:"created_at,omitempty"`
	UpdatedAt        string         `json:"updated_at,omitempty"`
}

type FirewallAPI struct {
	c           HTTPClient
	AuthToken   string
	Location    string
	ApiEndpoint string
	Firewall    *Firewall
	Firewalls   []Firewall
}

func (f *FirewallAPI) Init(c HTTPClient, authToken string, location string) error {
	f.c = c
	f.AuthToken = authToken
	f.Location = location
	f.ApiEndpoint = fmt.Sprintf(
		"https://api.idcloudhost.com/v1/%s/network/firewalls",
		f.Location,
	)
	req, err := http.NewRequest("GET", f.ApiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	r, err := f.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify endpoint: %v", err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("location: %s not found", f.Location)
	}
	return nil
}

// ListFirewalls lists all firewalls
func (f *FirewallAPI) ListFirewalls() error {
	req, err := http.NewRequest("GET", f.ApiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", f.AuthToken)
	r, err := f.c.Do(req)
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
	return json.NewDecoder(r.Body).Decode(&f.Firewalls)
}

// CreateFirewall creates a new firewall
func (f *FirewallAPI) CreateFirewall(firewall *Firewall) error {
	if err := validateFirewallRules(firewall.Rules); err != nil {
		return err
	}

	payloadJSON, err := json.Marshal(firewall)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", f.ApiEndpoint, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", f.AuthToken)

	r, err := f.c.Do(req)
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

	return json.NewDecoder(r.Body).Decode(&f.Firewall)
}

// UpdateFirewall updates an existing firewall
func (f *FirewallAPI) UpdateFirewall(uuid string, firewall *Firewall) error {
	if err := validateFirewallRules(firewall.Rules); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", f.ApiEndpoint, uuid)
	// For updates, only send rules
	payload := struct {
		Rules []FirewallRule `json:"rules"`
	}{
		Rules: firewall.Rules,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", f.AuthToken)

	r, err := f.c.Do(req)
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

	return json.NewDecoder(r.Body).Decode(&f.Firewall)
}

// DeleteFirewall deletes a firewall
func (f *FirewallAPI) DeleteFirewall(uuid string) error {
	url := fmt.Sprintf("%s/%s", f.ApiEndpoint, uuid)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", f.AuthToken)
	r, err := f.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if (r.StatusCode != http.StatusOK) && (r.StatusCode != http.StatusNoContent) {
		return fmt.Errorf("%v", r.StatusCode)
	}

	return nil
}

// AssignFirewall assigns a firewall to a virtual machine
func (f *FirewallAPI) AssignFirewall(firewallUUID string, vmUUID string) error {
	url := fmt.Sprintf("%s/%s/assign", f.ApiEndpoint, firewallUUID)

	payload := struct {
		VMUUID string `json:"vm_uuid"`
	}{
		VMUUID: vmUUID,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", f.AuthToken)

	r, err := f.c.Do(req)
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

// UnassignFirewall removes a firewall from a virtual machine
func (f *FirewallAPI) UnassignFirewall(firewallUUID string, vmUUID string) error {
	url := fmt.Sprintf("%s/%s/unassign", f.ApiEndpoint, firewallUUID)

	payload := struct {
		VMUUID string `json:"vm_uuid"`
	}{
		VMUUID: vmUUID,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", f.AuthToken)

	r, err := f.c.Do(req)
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

// validateFirewallRules validates the firewall rules
func validateFirewallRules(rules []FirewallRule) error {
	for _, rule := range rules {
		// Check direction (required)
		if rule.Direction != "inbound" && rule.Direction != "outbound" {
			return fmt.Errorf("firewall rule direction must be either 'inbound' or 'outbound'")
		}

		// Check protocol (required)
		if rule.Protocol != "tcp" && rule.Protocol != "udp" && rule.Protocol != "icmp" && rule.Protocol != "any" {
			return fmt.Errorf("firewall rule protocol must be one of 'tcp', 'udp', 'icmp', or 'any'")
		}

		// Check endpoint_spec_type (required)
		if rule.EndpointSpecType != "any" && rule.EndpointSpecType != "cidr" && rule.EndpointSpecType != "firewall" && rule.EndpointSpecType != "ip_prefixes" {
			return fmt.Errorf("firewall rule endpoint_spec_type must be one of 'any', 'cidr', 'firewall', or 'ip_prefixes'")
		}

		// Validate endpoint_spec is not nil (should be empty array if "any")
		if rule.EndpointSpec == nil {
			return fmt.Errorf("firewall rule endpoint_spec cannot be nil, use empty array for 'any' type")
		}
	}
	return nil
}
