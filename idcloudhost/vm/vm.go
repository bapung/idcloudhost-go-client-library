package vm

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bapung/idcloudhost-go-client-library/idcloudhost/disk"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VirtualMachineAPI struct {
	c              HTTPClient
	AuthToken      string
	Location       string
	BillingAccount int
	ApiEndpoint    string
	VM             VM
	VMList         []VM
	Replica        Replica
	ReplicaList    []Replica
	ResourcePools  []ResourcePool
}

type VM struct {
	Backup         bool               `json:"backup"`
	BillingAccount int                `json:"billing_account"`
	CreatedAt      string             `json:"created_at"`
	Description    string             `json:"description"`
	Hostname       string             `json:"hostname"`
	HypervisorId   string             `json:"hypervisor_id"`
	Id             int                `json:"id"`
	MACAddress     string             `json:"mac"`
	Memory         int                `json:"memory"`
	Name           string             `json:"name"`
	OSName         string             `json:"os_name"`
	OSVersion      string             `json:"os_version"`
	PrivateIPv4    string             `json:"private_ipv4"`
	Status         string             `json:"status"`
	Storage        []disk.DiskStorage `json:"storage"`
	Tags           []string           `json:"tags"`
	UpdatedAt      string             `json:"updated_at"`
	UserId         int                `json:"user_id"`
	Username       string             `json:"username"`
	UUID           string             `json:"uuid"`
	VCPU           int                `json:"vcpu"`
}

type NewVM struct {
	Backup          bool   `validate:"-" default:"false"`
	BillingAccount  int    `validate:"-" default:"0"`
	Description     string `validate:"-"`
	Disks           int    `validate:"required|int|min:20|max:240"`
	Username        string `validate:"validateUsername"`
	InitialPassword string `validate:"required|validatePassword"`
	OSName          string `validate:"required|validateOSName"`
	OSVersion       string `validate:"required|validateOSVersion"`
	PublicKey       string `validate:"-"`
	Name            string `validate:"required|validateName"`
	Memory          int    `validate:"required|int|min:1024|max:65536"`
	SourceReplica   string `validate:"-"`
	SourceUUID      string `validate:"-"`
	VCPU            int    `validate:"required|int|min:1|max:16"`
	ReservePublicIP bool   `validate:"-" default:"true"`
}

type ResourcePool struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CPUCount    int    `json:"cpu_count"`
	RAMSize     int    `json:"ram_size"`
	DiskSize    int    `json:"disk_size"`
}

type Replica struct {
	ID        int    `json:"id,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	Name      string `json:"name"`
	VMUUID    string `json:"vm_uuid,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
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

func (vm *VirtualMachineAPI) Create(v NewVM) error {
	if err := validateVmCreateFields(&v); err != nil {
		return err
	}
	data := url.Values{}
	data.Set("billing_account_id", strconv.Itoa(v.BillingAccount))
	data.Set("disks", strconv.Itoa(v.Disks))
	data.Set("name", v.Name)
	data.Set("username", v.Username)
	data.Set("password", v.InitialPassword)
	data.Set("os_name", v.OSName)
	data.Set("os_version", v.OSVersion)
	data.Set("vcpu", strconv.Itoa(v.VCPU))
	data.Set("ram", strconv.Itoa(v.Memory))
	data.Set("backup", strconv.FormatBool(v.Backup))
	data.Set("reserve_public_ip", strconv.FormatBool(v.ReservePublicIP))
	if v.Description != "" {
		data.Set("description", v.Description)
	}
	if v.PublicKey != "" {
		data.Set("public_key", v.PublicKey)
	}
	if v.SourceReplica != "" {
		data.Set("source_replica", v.SourceReplica)
	}
	if v.SourceUUID != "" {
		data.Set("source_uuid", v.SourceUUID)
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	return json.NewDecoder(r.Body).Decode(&vm.VM)
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	return json.NewDecoder(r.Body).Decode(&vm.VM)
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	return json.Unmarshal(body, &vm.VMList)
}

func (vm *VirtualMachineAPI) Modify(v VM) error {
	if err := validateVmModifyFields(&v); err != nil {
		return err
	}
	data := url.Values{}
	if v.Name == vm.VM.Name && v.VCPU == vm.VM.VCPU && v.Memory == vm.VM.Memory {
		return fmt.Errorf("name or VCPU or RAM value does not changed, not updating")
	}
	data.Set("uuid", v.UUID)
	data.Set("name", v.Name)
	// workaround to idcloudhost API bug, cannot be set if none is changed.
	if v.Memory != vm.VM.Memory || v.VCPU != vm.VM.VCPU {
		data.Set("ram", strconv.Itoa(v.Memory))
		data.Set("vcpu", strconv.Itoa(v.VCPU))
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

func (vm *VirtualMachineAPI) Clone(uuid string, cloneName string) error {
	backupApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "clone")
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("name", cloneName)
	req, err := http.NewRequest("POST", backupApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) ToggleAutoBackup(uuid string) error {
	backupApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "backup")
	data := url.Values{}
	data.Set("uuid", uuid)
	req, err := http.NewRequest("POST", backupApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) ReleasePublicIP(uuid string) error {
	releaseIPApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "release-ip")
	data := url.Values{}
	data.Set("uuid", uuid)
	req, err := http.NewRequest("POST", releaseIPApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) ReservePublicIP(uuid string) error {
	reserveIPApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "reserve-ip")
	data := url.Values{}
	data.Set("uuid", uuid)
	req, err := http.NewRequest("POST", reserveIPApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) RebuildFromReplica(uuid string, replicaUUID string) error {
	rebuildApiEndpoint := fmt.Sprintf("%s/%s/rebuild", vm.ApiEndpoint, "replica")
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("replica_uuid", replicaUUID)
	req, err := http.NewRequest("POST", rebuildApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) DeleteReplica(uuid string, replicaUUID string) error {
	replicaApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "replica")
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("replica_uuid", replicaUUID)
	req, err := http.NewRequest("DELETE", replicaApiEndpoint,
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

func (vm *VirtualMachineAPI) ListReplicas(uuid string) error {
	replicaApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "replica")
	url := fmt.Sprintf("%s?uuid=%s", replicaApiEndpoint, uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := vm.c.Do(req)
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
	return json.NewDecoder(r.Body).Decode(&vm.ReplicaList)
}

func (vm *VirtualMachineAPI) CreateReplica(uuid string, name string) error {
	replicaApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "replica")
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("name", name)
	req, err := http.NewRequest("POST", replicaApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.Replica); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) StartVM(uuid string) error {
	startApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "start")
	data := url.Values{}
	data.Set("uuid", uuid)
	req, err := http.NewRequest("POST", startApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) StopVM(uuid string) error {
	stopApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "stop")
	data := url.Values{}
	data.Set("uuid", uuid)
	req, err := http.NewRequest("POST", stopApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) BootMedia(uuid string, mediaType string) error {
	bootApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "boot")
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("type", mediaType)
	req, err := http.NewRequest("POST", bootApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) ChangePassword(uuid string, newPassword string) error {
	passwordApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "password")
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("password", newPassword)
	req, err := http.NewRequest("POST", passwordApiEndpoint,
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

func (vm *VirtualMachineAPI) Reinstall(uuid string, osName string, osVersion string, username string, password string) error {
	reinstallApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "reinstall")
	data := url.Values{}
	data.Set("uuid", uuid)
	data.Set("os_name", osName)
	data.Set("os_version", osVersion)
	data.Set("username", username)
	data.Set("password", password)
	req, err := http.NewRequest("POST", reinstallApiEndpoint,
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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", r.StatusCode)
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VM); err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachineAPI) ListResourcePools() error {
	url := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "resource-pools")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := vm.c.Do(req)
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
	return json.NewDecoder(r.Body).Decode(&vm.ResourcePools)
}
