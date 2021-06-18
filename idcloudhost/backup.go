package idcloudhost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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
