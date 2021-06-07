package idcloudhost

import (
	"fmt"
	"net/http"
	"net/url"
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

func (vm *VirtualMachineAPI) AttachDisk(v map[string]interface{}) {
	var c HTTPClient = &http.Client{}
	var disk *DiskStorage
	diskApiEndpoint := fmt.Sprintf("%s/%s", vm.ApiEndpoint, "disk")
	data := url.Values{}
	data.Set("uuid", v["uuid"].(string))
	data.Set("size_gb", strconv.Iota(v["disk_size"].(int))))
	req, err := http.NewRequest("POST", diskApiEndpoint,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", vm.AuthToken)
	r, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(&disk)
}

