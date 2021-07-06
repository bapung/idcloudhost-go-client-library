package idcloudhost

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type DiskAPI struct {
	c           HTTPClient
	AuthToken   string
	Location    string
	ApiEndpoint string
	vmUUID      string
	DiskList    *[]DiskStorage
	Disk        *DiskStorage
}

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

func DiskVmNotSpecifiedError() error {
	return fmt.Errorf("disk API must be called after VM UUID is specified with Bind()")
}

func DiskNotFoundError() error {
	return fmt.Errorf("specified disk not found")
}

func (d *DiskAPI) Init(c HTTPClient, authToken string, location string) error {
	d.c = c
	d.AuthToken = authToken
	d.Location = location
	d.ApiEndpoint = fmt.Sprintf(
		"https://api.idcloudhost.com/v1/%s/user-resource/vm/storage",
		d.Location,
	)
	r, err := http.Get(d.ApiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("location: %s not found", d.Location)
	}
	return nil
}

func (d *DiskAPI) Bind(vmUUID string) {
	d.vmUUID = vmUUID
}

func (d *DiskAPI) Get(uuid string) error {
	if d.vmUUID == "" {
		return DiskVmNotSpecifiedError()
	}
	var v VirtualMachineAPI
	v.Init(d.c, d.AuthToken, d.Location)
	err := v.Get(d.vmUUID)
	if err != nil {
		return err
	}
	d.DiskList = &v.VM.Storage

	for _, disk := range *(d.DiskList) {
		if disk.UUID == uuid {
			d.Disk = &disk
			return nil
		}
	}
	return DiskNotFoundError()
}

func (d *DiskAPI) Create(diskSize int) error {
	if d.vmUUID == "" {
		return DiskVmNotSpecifiedError()
	}
	if err := validateDisks(diskSize); err != nil {
		return err
	}
	data := url.Values{}
	data.Set("uuid", d.vmUUID)
	data.Set("size_gb", strconv.Itoa(diskSize))
	req, err := http.NewRequest("POST", d.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", d.AuthToken)
	r, err := d.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()

	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&d.Disk)
}

func (d *DiskAPI) Delete(diskUUID string) error {
	if d.vmUUID == "" {
		return DiskVmNotSpecifiedError()
	}
	var resp map[string]interface{}
	data := url.Values{}
	data.Set("uuid", d.vmUUID)
	data.Set("storage_uuid", diskUUID)
	req, err := http.NewRequest("DELETE", d.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", d.AuthToken)
	r, err := d.c.Do(req)
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
	if !resp["success"].(bool) {
		return fmt.Errorf("unknown error")
	}
	return nil
}

func (d *DiskAPI) Modify(diskUUID string, newDiskSize int) error {
	if d.vmUUID == "" {
		return DiskVmNotSpecifiedError()
	}
	data := url.Values{}
	data.Set("uuid", d.vmUUID)
	data.Set("disk_uuid", diskUUID)
	data.Set("size_gb", strconv.Itoa(newDiskSize))
	req, err := http.NewRequest("PATCH", d.ApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", d.AuthToken)
	r, err := d.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	if err = json.NewDecoder(r.Body).Decode(&d.Disk); err != nil {
		return err
	}
	return nil
}
