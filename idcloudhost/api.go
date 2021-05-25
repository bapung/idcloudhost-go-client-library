package idcloudhost

import (
	"net/http"
)

type API interface {
	Init(string, string) error
	Get(string) error
	Create(map[string]interface{}) error
	Modify(map[string]interface{}) error
	Delete(string) error
}

type APIClient struct {
	APIs []API
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(authToken string, loc string) (APIClient, error) {
	vmApi := VirtualMachineAPI{}
	// add more client here
	client := []API{&vmApi}
	for _, c := range client {
		c.Init(authToken, loc)
	}
	return APIClient{client}, nil
}
