package idcloudhost

import (
	"net/http"
)

type APIClient struct {
	VM   *VirtualMachineAPI
	Disk *DiskAPI
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(authToken string, loc string) (*APIClient, error) {
	c := http.Client{}
	var ApiClient = APIClient{
		VM:   &VirtualMachineAPI{},
		Disk: &DiskAPI{},
	}

	ApiClient.VM.Init(&c, authToken, loc)

	return &ApiClient, nil
}
