package idcloudhost

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type VPCNetworkAPI struct {
	c              HTTPClient
	AuthToken      string
	Location       string
	ApiEndpoint    string
	VPCNetworkList *[]VPCNetwork
	VPCNetwork     *VPCNetwork
}

type VPCNetwork struct {
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	UUID           string   `json:"uuid"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	VlanId         int      `json:"vlan_id"`
	IsDefault      bool     `json:"is_default"`
	VMUUIDList     []string `json:"vm_uuids"`
	Subnet         string   `json:"subnet"`
	ResourcesCount int      `json:"resources_count"`
}

func (vpc *VPCNetworkAPI) Init(c HTTPClient, authToken string, location string) error {
	vpc.c = c
	vpc.AuthToken = authToken
	vpc.Location = location
	vpc.ApiEndpoint = fmt.Sprintf(
		"https://api.idcloudhost.com/v1/%s/network/network",
		vpc.Location,
	)
	r, err := http.Get(vpc.ApiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("location: %s not found", vpc.Location)
	}
	return nil
}

func (vpc *VPCNetworkAPI) List() error {
	url := fmt.Sprintf("%ss", vpc.ApiEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vpc.AuthToken)
	r, err := vpc.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&vpc.VPCNetworkList)
}

func (vpc *VPCNetworkAPI) Get(UUID string) error {
	var url = fmt.Sprintf("%s/%s", vpc.ApiEndpoint, UUID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vpc.AuthToken)
	r, err := vpc.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&vpc.VPCNetwork)
}

func (vpc *VPCNetworkAPI) Create(name string) error {
	data := url.Values{}
	data.Set("name", name)
	req, err := http.NewRequest("POST", vpc.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("apiKey", vpc.AuthToken)
	r, err := vpc.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&vpc.VPCNetwork)
}

func (vpc *VPCNetworkAPI) SetDefault(UUID string) error {
	var url = fmt.Sprintf("%s/%s/default", vpc.ApiEndpoint, UUID)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vpc.AuthToken)
	r, err := vpc.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&vpc.VPCNetwork)
}

func (vpc *VPCNetworkAPI) Delete(UUID string) error {
	var url = fmt.Sprintf("%s/%s", vpc.ApiEndpoint, UUID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vpc.AuthToken)
	r, err := vpc.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return nil
}
