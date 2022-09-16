package idcloudhost

import (
	"net/http"
)

type APIClient struct {
	VM         		*VirtualMachineAPI
	Disk       		*DiskAPI
	FloatingIP 		*FloatingIPAPI
	LoadBalancer 	*LoadBalancerAPI
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(authToken string, loc string) (*APIClient, error) {
	c := http.Client{}
	var ApiClient = APIClient{
		VM:         &VirtualMachineAPI{},
		Disk:       &DiskAPI{},
		FloatingIP: &FloatingIPAPI{},
		LoadBalancer: &LoadBalancerAPI{},
	}

	ApiClient.VM.Init(&c, authToken, loc)
	ApiClient.Disk.Init(&c, authToken, loc)
	ApiClient.FloatingIP.Init(&c, authToken, loc)
	ApiClient.LoadBalancer.Init(&c, authToken, loc)

	return &ApiClient, nil
}
