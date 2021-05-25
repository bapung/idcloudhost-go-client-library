package idcloudhost

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type VirtualMachineAPI struct {
	AuthToken   string
	Location    string
	ApiEndpoint string
	VM          *VM
	VMMap       map[string]interface{}
	VMList      []VM
	VMListMap   []map[string]interface{}
}

type vmList struct {
	vm []VM
}

type VM struct {
	Backup         bool           `json:"backup,omitempty"`
	BillingAccount int            `json:"billing_account,omitempty"`
	CreatedAt      string         `json:"created_at,omitempty"`
	Description    string         `json:"description"`
	Hostname       string         `json:"hostname,omitempty"`
	HypervisorId   string         `json:"hypervisor_id,omitempty"`
	Id             int            `json:"id,omitempty"`
	MACAddress     string         `json:"mac,omitempty"`
	MemoryM        int            `json:"memory"`
	Name           string         `json:"name"`
	OSName         string         `json:"os_name"`
	OSVersion      string         `json:"os_version"`
	PrivateIPv4    string         `json:"private_ipv4,omitempty"`
	Status         string         `json:"status,omitempty"`
	Storage        *[]DiskStorage `json:"storage,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	UpdatedAt      string         `json:"updated_at,omitempty"`
	UserId         int            `json:"user_id,omitempty"`
	Username       string         `json:"username"`
	UUID           string         `json:"uuid,omitempty"`
	VCPU           int            `json:"vcpu"`
}

type omit *struct{}

type NewVM struct {
	Backup          string `json:"backup",default:"false"`
	BillingAccount  int    `json:"billing_account_id",default:0`
	Description     string `json:"description"`
	Disks           int    `json:"disks"`
	InitialPassword string `json:"password"`
	OSName          string `json:"os_name"`
	OSVersion       string `json:"os_version"`
	PublicKey       string `json:"public_key,omitempty"`
	MemoryM         string `json:"ram"`
	SourceReplica   string `json:"source_replica,omitempty"`
	SourceUUID      string `json:"source_uuid,omitempty"`
	VCPU            int    `json:"vcpu"`
	PublicIP        string `json:"reserve_public_ip",default:"true"`
}

func (vm *VirtualMachineAPI) Init(authToken string, location string) error {
	vm.AuthToken = authToken
	vm.Location = location
	vm.ApiEndpoint = fmt.Sprintf(
		"https://api.idcloudhost.com/v1/%s/user-resource/vm",
		vm.Location,
	)
	r, err := http.Get(vm.ApiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == http.StatusNotFound {
		return errors.New(fmt.Sprintf("Location: %s not found", vm.Location))
	}
	return nil
}

func (vm *VirtualMachineAPI) Create(newVm map[string]interface{}) error {
	var c HTTPClient
	c = &http.Client{}
	data := url.Values{}
	data.Set("backup", v["backup"])
	data.Set("billing_account_id", v["billing_account"].(string))
	data.Set("description", v["description"])
	data.Set("disks", v["disks"].(string))
	data.Set("password", v["password"])
	data.Set("os_name", v["os_name"])
	data.Set("os_version", v["os_version"])
	data.Set("vcpu", v["vcpu"].(string))
	data.Set("ram", v["ram"].(string))
	if v["public_key"] != "" {
		data.Set("public_key", v["public_key"])
	}
	if v["source_replica"] != "" {
		data.Set("source_replica", v["source_replica"])
	}
	if v["source_uuid "] != "" {
		data.Set("source_uuid", v["source_uuid)"])
	}
	if v["reserve_public_ip"] != "" {
		data.Set("reserve_public_ip", v["reserve_public_ip"].(string))
	}
	req, err := http.NewRequest("POST", vm.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&vm.VMMap)
}

func (vm *VirtualMachineAPI) Get(uuid string) error {
	var c HTTPClient
	c = &http.Client{}
	url := fmt.Sprintf("%s?uuid=%s", vm.ApiEndpoint, uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&vm.VMMap)
}

func (vm *VirtualMachineAPI) ListAll() error {
	var c HTTPClient
	c = &http.Client{}
	url := fmt.Sprintf("%s/list", vm.ApiEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	bodyByte, err := ioutil.ReadAll(r.Body)
	return json.Unmarshal(bodyByte, &vm.VMListMap)
}

func (vm *VirtualMachineAPI) Modify(v map[string]interface{}) error {
	var c HTTPClient
	c = &http.Client{}
	data := url.Values{}
	data.Set("uuid", v["uuid"].(string))
	data.Set("name", v["name"].(string))
	data.Set("ram", v["ram"].(string))
	data.Set("vcpu", v["vcpu"].(string))
	req, err := http.NewRequest("PATCH", vm.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer r.Body.Close()
	return checkError(r.StatusCode)
}

func (vm *VirtualMachineAPI) Delete(uuid string) error {
	var c HTTPClient
	c = &http.Client{}
	data := url.Values{}
	data.Set("uuid", uuid)
	req, err := http.NewRequest("DELETE", vm.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer r.Body.Close()
	return checkError(r.StatusCode)
}
