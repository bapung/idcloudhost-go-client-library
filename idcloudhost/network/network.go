package network

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

type Network struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	UUID        string `json:"uuid,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
	Default     bool   `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type NetworkAPI struct {
	c           HTTPClient
	AuthToken   string
	Location    string
	ApiEndpoint string
	Network     *Network
	Networks    []Network
}

func (n *NetworkAPI) Init(c HTTPClient, authToken string, location string) error {
	n.c = c
	n.AuthToken = authToken
	n.Location = location
	n.ApiEndpoint = fmt.Sprintf(
		"https://api.idcloudhost.com/v1/%s/network/private_networks",
		n.Location,
	)
	r, err := http.Get(n.ApiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("location: %s not found", n.Location)
	}
	return nil
}

// GetNetwork gets details for a specific network
func (n *NetworkAPI) GetNetwork(uuid string) error {
	url := fmt.Sprintf("%s/%s", n.ApiEndpoint, uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", n.AuthToken)
	r, err := n.c.Do(req)
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
	return json.NewDecoder(r.Body).Decode(&n.Network)
}

// ListNetworks lists all private networks
func (n *NetworkAPI) ListNetworks() error {
	req, err := http.NewRequest("GET", n.ApiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", n.AuthToken)
	r, err := n.c.Do(req)
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
	return json.NewDecoder(r.Body).Decode(&n.Networks)
}

// CreateDefaultNetwork creates or retrieves the default network
func (n *NetworkAPI) CreateDefaultNetwork() error {
	req, err := http.NewRequest("POST", n.ApiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", n.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	r, err := n.c.Do(req)
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
	return json.NewDecoder(r.Body).Decode(&n.Network)
}

// DeleteNetwork deletes a specific network
func (n *NetworkAPI) DeleteNetwork(uuid string) error {
	url := fmt.Sprintf("%s/%s", n.ApiEndpoint, uuid)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", n.AuthToken)
	r, err := n.c.Do(req)
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

// SetAsDefault sets a network as the default network
func (n *NetworkAPI) SetAsDefault(uuid string) error {
	url := fmt.Sprintf("%s/%s/default", n.ApiEndpoint, uuid)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", n.AuthToken)
	r, err := n.c.Do(req)
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
	return json.NewDecoder(r.Body).Decode(&n.Network)
}

// UpdateNetwork updates a network's name
func (n *NetworkAPI) UpdateNetwork(uuid string, name string) error {
	url := fmt.Sprintf("%s/%s", n.ApiEndpoint, uuid)

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
	req.Header.Set("apiKey", n.AuthToken)

	r, err := n.c.Do(req)
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

	return json.NewDecoder(r.Body).Decode(&n.Network)
}
