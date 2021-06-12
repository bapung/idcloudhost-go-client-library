package idcloudhost

import (
	"net/http"
)

type API interface {
	Init(HTTPClient, string, string) error
	Get(string) error
	Create(map[string]interface{}) error
	Modify(map[string]interface{}) error
	Delete(string) error
}

type APIClient struct {
	APIs map[string]API
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(authToken string, loc string) (*APIClient, error) {
	api := APIClient{}
	c := http.Client{}
	var m = map[string]API{
		"vm": &VirtualMachineAPI{},
	}
	api.APIs = m
	// add more client here
	for _, k := range api.APIs {
		k.Init(&c, authToken, loc)
	}
	return &api, nil
}
