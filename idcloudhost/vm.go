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

func (vm *VirtualMachineAPI) Create(newVm *NewVM) error {
	var c HTTPClient
	c = &http.Client{}
	data := url.Values{}
	data.Set("backup", newVm.Backup)
	data.Set("billing_account_id", fmt.Sprint(newVm.BillingAccount))
	data.Set("description", newVm.Description)
	data.Set("disks", fmt.Sprint(newVm.Disks))
	data.Set("password", newVm.InitialPassword)
	data.Set("os_name", newVm.OSName)
	data.Set("os_version", newVm.OSVersion)
	data.Set("vcpu", fmt.Sprint(newVm.VCPU))
	data.Set("ram", string(newVm.MemoryM))
	if newVm.PublicIP != "" {
		data.Set("public_key", newVm.PublicIP)
	}
	if newVm.SourceReplica != "" {
		data.Set("source_replica", newVm.SourceReplica)
	}
	if newVm.SourceUUID != "" {
		data.Set("source_uuid", newVm.SourceUUID)
	}
	if newVm.SourceUUID != "" {
		data.Set("reserve_public_ip", newVm.PublicIP)
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
	return json.Unmarshal(bodyByte, &vm.VMList)
}

func (vm *VirtualMachineAPI) Modify(uuid string, name string, ram int, vcpu int) error {
	var c HTTPClient
	c = &http.Client{}
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("name", name)
	data.Set("ram", fmt.Sprint(ram))
	data.Set("vcpu", fmt.Sprint(vcpu))
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
