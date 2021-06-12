package idcloudhost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type VirtualMachineAPI struct {
	c              HTTPClient
	AuthToken      string
	Location       string
	BillingAccount int
	ApiEndpoint    string
	VM             *VM
	VMMap          map[string]interface{}
	VMList         []VM
	VMListMap      []map[string]interface{}
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
	MemoryM        int            `json:"ram"`
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

func (vm *VirtualMachineAPI) Init(c HTTPClient, authToken string, location string) error {
	vm.c = c
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
		return fmt.Errorf("location: %s not found", vm.Location)
	}
	return nil
}

func (vm *VirtualMachineAPI) Create(v map[string]interface{}) error {
	if err := validateVirtualMachineParam(v); err != nil {
		return err
	}
	data := url.Values{}
	data.Set("backup", strconv.FormatBool(v["backup"].(bool)))
	data.Set("billing_account_id", strconv.Itoa(v["billing_account"].(int)))
	data.Set("disks", strconv.Itoa(v["disks"].(int)))
	data.Set("name", v["name"].(string))
	data.Set("username", v["username"].(string))
	data.Set("password", v["password"].(string))
	data.Set("os_name", v["os_name"].(string))
	data.Set("os_version", v["os_version"].(string))
	data.Set("vcpu", strconv.Itoa(v["vcpu"].(int)))
	data.Set("ram", strconv.Itoa(v["ram"].(int)))
	if v["description"] != nil {
		data.Set("description", v["description"].(string))
	}
	if v["public_key"] != nil {
		data.Set("public_key", v["public_key"].(string))
	}
	if v["source_replica"] != nil {
		data.Set("source_replica", v["source_replica"].(string))
	}
	if v["source_uuid "] != nil {
		data.Set("source_uuid", v["source_uuid)"].(string))
	}
	if v["reserve_public_ip"] != nil {
		data.Set("reserve_public_ip", strconv.FormatBool(v["reserve_public_ip"].(bool)))
	}
	req, err := http.NewRequest("POST", vm.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := vm.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&vm.VMMap)
}

func (vm *VirtualMachineAPI) Get(uuid string) error {
	url := fmt.Sprintf("%s?uuid=%s", vm.ApiEndpoint, uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := vm.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&vm.VMMap)
}

func (vm *VirtualMachineAPI) ListAll() error {
	url := fmt.Sprintf("%s/list", vm.ApiEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := vm.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	bodyByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	return json.Unmarshal(bodyByte, &vm.VMListMap)
}

func (vm *VirtualMachineAPI) Modify(v map[string]interface{}) error {
	data := url.Values{}
	data.Set("uuid", v["uuid"].(string))
	if v["name"] != nil {
		data.Set("name", v["name"].(string))
	}
	if v["ram"] != nil {
		data.Set("ram", strconv.Itoa(v["ram"].(int)))
	}
	if v["vcpu"] != nil {
		data.Set("vcpu", strconv.Itoa(v["vcpu"].(int)))
	}
	req, err := http.NewRequest("PATCH", vm.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r, err := vm.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	return checkError(r.StatusCode)
}

func (vm *VirtualMachineAPI) Delete(uuid string) error {
	data := url.Values{}
	data.Set("uuid", uuid)
	req, err := http.NewRequest("DELETE", vm.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r, err := vm.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	return checkError(r.StatusCode)
}
