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
	Type        string `json:"type"`
	Protocol    string `json:"protocol"`
	PortRange   string `json:"port_range,omitempty"`
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
	Description string `json:"description,omitempty"`
}

type Firewall struct {
	ID          int            `json:"id,omitempty"`
	Name        string         `json:"name"`
	UUID        string         `json:"uuid,omitempty"`
	UserID      int            `json:"user_id,omitempty"`
	Description string         `json:"description,omitempty"`
	Rules       []FirewallRule `json:"rules"`
	CreatedAt   string         `json:"created_at,omitempty"`
	UpdatedAt   string         `json:"updated_at,omitempty"`
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
	payloadJSON, err := json.Marshal(firewall)
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

	if r.StatusCode != http.StatusOK {
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
		// Check type (required)
		if rule.Type != "ingress" && rule.Type != "egress" {
			return fmt.Errorf("firewall rule type must be either 'ingress' or 'egress'")
		}

		// Check protocol (required)
		if rule.Protocol != "tcp" && rule.Protocol != "udp" && rule.Protocol != "icmp" {
			return fmt.Errorf("firewall rule protocol must be one of 'tcp', 'udp', or 'icmp'")
		}

		// Check port range format for tcp/udp
		// TODO: add more detailed port range validation if needed
		// For now, we accept any non-empty port range for non-ICMP protocols
	}
	return nil
}
