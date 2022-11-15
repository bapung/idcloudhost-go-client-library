package api

import (
	"net/http"

	"github.com/bapung/idcloudhost-go-client-library/idcloudhost/disk"
	"github.com/bapung/idcloudhost-go-client-library/idcloudhost/floatingip"
	"github.com/bapung/idcloudhost-go-client-library/idcloudhost/loadbalancer"
	"github.com/bapung/idcloudhost-go-client-library/idcloudhost/vm"
)

type APIClient struct {
	VM         *vm.VirtualMachineAPI
	Disk       *disk.DiskAPI
	FloatingIP *floatingip.FloatingIPAPI
	LB         *loadbalancer.LoadBalancerAPI
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
		LB:         &loadbalancer.LoadBalancerAPI{},
	}

	ApiClient.VM.Init(&c, authToken, loc)
	ApiClient.Disk.Init(&c, authToken, loc)
	ApiClient.FloatingIP.Init(&c, authToken, loc)
	ApiClient.LB.Init(&c, authToken, loc)

	return &ApiClient, nil
}
