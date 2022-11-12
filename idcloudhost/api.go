package idcloudhost

import (
	"net/http"

	"github.com/bapung/idcloudhost-go-client-library/disk"
	"github.com/bapung/idcloudhost-go-client-library/floatingip"
	"github.com/bapung/idcloudhost-go-client-library/vm"
)

type APIClient struct {
	VM         *vm.VirtualMachineAPI
	Disk       *disk.DiskAPI
	FloatingIP *floatingip.FloatingIPAPI
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(authToken string, loc string) (*APIClient, error) {
	c := http.Client{}
	var ApiClient = APIClient{
		VM:         &vm.VirtualMachineAPI{},
		Disk:       &disk.DiskAPI{},
		FloatingIP: &floatingip.FloatingIPAPI{},
	}

	ApiClient.VM.Init(&c, authToken, loc)
	ApiClient.Disk.Init(&c, authToken, loc)
	ApiClient.FloatingIP.Init(&c, authToken, loc)

	return &ApiClient, nil
}
