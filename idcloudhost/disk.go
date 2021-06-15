package idcloudhost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type DiskStorage struct {
	CreatedAt string   `json:"created_at"`
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Pool      string   `json:"pool"`
	Primary   bool     `json:"primary"`
	Replica   []string `json:"replica"`
	Shared    bool     `json:"shared"`
	SizeGB    int      `json:"size"`
	Type      string   `json:"type"`
	UpdatedAt string   `json:"updated_at"`
	UserId    int      `json:"user_id"`
	UUID      string   `json:"uuid"`
}

func (vm *VirtualMachineAPI) AttachDisk(vmUUID string, diskSize int) (*DiskStorage, error) {
	var disk *DiskStorage
	diskApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "disk")
	data := url.Values{}
	data.Set("uuid", vmUUID)
	data.Set("size_gb", strconv.Itoa(diskSize))
	req, err := http.NewRequest("POST", diskApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := vm.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return nil, err
	}
	if err := json.NewDecoder(r.Body).Decode(&disk); err != nil {
		return nil, err
	}
	return disk, nil
}

func (vm *VirtualMachineAPI) DeleteDisk(vmUUID string, diskUUID string) error {
	var resp map[string]interface{}
	diskApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "disk")
	data := url.Values{}
	data.Set("uuid", vmUUID)
	data.Set("storage_uuid", diskUUID)
	req, err := http.NewRequest("DELETE", diskApiEndpoint,
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
	if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}
	if resp["success"].(bool) {
		return fmt.Errorf("unknown error")
	}
	return nil
}

func (vm *VirtualMachineAPI) ResizeDisk(vmUUID string, diskUUID string, newDiskSize int) error {
	var resp map[string]interface{}
	diskApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "disk")
	data := url.Values{}
	data.Set("uuid", vmUUID)
	data.Set("disk_uuid", diskUUID)
	data.Set("size_gb", strconv.Itoa(newDiskSize))
	req, err := http.NewRequest("DELETE", diskApiEndpoint,
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
	if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}
	if resp["success"].(bool) {
		return fmt.Errorf("unknown error")
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
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	if err = json.NewDecoder(r.Body).Decode(&vm.VMMap); err != nil {
		return err
	}
	return nil
}
