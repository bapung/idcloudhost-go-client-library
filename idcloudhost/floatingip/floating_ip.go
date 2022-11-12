package floatingip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"errors"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type FloatingIP struct {
	ID             int    `json:"id,omitempty"`
	Address        string `json:"address,omitempty"`
	UserID         int    `json:"user_id,omitempty"`
	BillingAccount int    `json:"billing_account_id"`
	Type           string `json:"type,omitempty"`
	NetworkID      string `json:"network_id,omitempty"`
	Name           string `json:"name"`
	Enabled        bool   `json:"enabled,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	AssignedTo     string `json:"assigned_to,omitempty"`
}

type FloatingIPAPI struct {
	c           HTTPClient
	AuthToken   string
	Location    string
	ApiEndpoint string
	FloatingIP  *FloatingIP
}

func (ip *FloatingIPAPI) Init(c HTTPClient, authToken string, location string) error {
	ip.c = c
	ip.AuthToken = authToken
	ip.Location = location
	ip.ApiEndpoint = fmt.Sprintf(
		"https://api.idcloudhost.com/v1/%s/network/ip_addresses",
		ip.Location,
	)
	r, err := http.Get(ip.ApiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("location: %s not found", ip.Location)
	}
	return nil
}

func (ip *FloatingIPAPI) Create(name string, billingAccountId int) error {
	var payloadJSON, err = json.Marshal(
		&FloatingIP{
			Name:           name,
			BillingAccount: billingAccountId,
		})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", ip.ApiEndpoint,
		bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("apiKey", ip.AuthToken)
	r, err := ip.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v",r.StatusCode))
	}
	return json.NewDecoder(r.Body).Decode(&ip.FloatingIP)
}

func (ip *FloatingIPAPI) Get(IPAddress string) error {
	var url = fmt.Sprintf("%s/%s", ip.ApiEndpoint, IPAddress)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", ip.AuthToken)
	r, err := ip.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v",r.StatusCode))
	}
	return json.NewDecoder(r.Body).Decode(&ip.FloatingIP)
}

func (ip *FloatingIPAPI) Update(name string, billingAccountId int, targetIPAddress string) error {
	var url = fmt.Sprintf("%s/%s", ip.ApiEndpoint, targetIPAddress)
	var payloadJSON, err = json.Marshal(
		&FloatingIP{
			Name:           name,
			BillingAccount: billingAccountId,
		})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", url,
		bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("apiKey", ip.AuthToken)
	r, err := ip.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v",r.StatusCode))
	}
	return json.NewDecoder(r.Body).Decode(&ip.FloatingIP)
}

func (ip *FloatingIPAPI) Delete(IPAddress string) error {
	var url = fmt.Sprintf("%s/%s", ip.ApiEndpoint, IPAddress)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", ip.AuthToken)
	r, err := ip.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v",r.StatusCode))
	}
	return nil
}

func (ip *FloatingIPAPI) Assign(IPAddress string, targetVMUUID string) error {
	var url = fmt.Sprintf("%s/%s/assign", ip.ApiEndpoint, IPAddress)
	var payloadJSON = []byte(
		fmt.Sprintf("{ \"vm_uuid\": \"%s\" }", targetVMUUID))
	req, err := http.NewRequest("POST", url,
		bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("apiKey", ip.AuthToken)
	r, err := ip.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v",r.StatusCode))
	}
	return json.NewDecoder(r.Body).Decode(&ip.FloatingIP)
}

func (ip *FloatingIPAPI) Unassign(IPAddress string) error {
	var url = fmt.Sprintf("%s/%s/unassign", ip.ApiEndpoint, IPAddress)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("apiKey", ip.AuthToken)
	r, err := ip.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v",r.StatusCode))
	}
	return json.NewDecoder(r.Body).Decode(&ip.FloatingIP)
}
