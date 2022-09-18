package idcloudhost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
)

type LoadBalancerAPI struct {
	c           		HTTPClient
	AuthToken   		string
	Location    		string
	ApiEndpoint 		string
	LoadBalancerList    *[]LoadBalancer
}


type LoadBalancer struct  {
	UUID 				string 				`json:"uuid"`
	NetworkUUID 		string 				`json:"network_uuid"`
	DisplayName			string				`json:"display_name,omitempty"`
	UserUID 			int 				`json:"user_id"`
	BillingAccount 		int 				`json:"billing_account_id"`
	CreatedAt			string  			`json:"created_at"`
	UpdatedAt			string  			`json:"updated_at"`
	IsDeleted			bool    			`json:"is_deleted"`
	PrivateIPv4     	string  			`json:"private_address"`
	ReservePublicIP 	string				`json:"reserve_public_ip,omitempty"`
	ForwardingRules 	[]ForwardingRule 	`json:"forwarding_rules"`
	Targets				[]ForwardingTarget 	`json:"targets"`
}

type ForwardingRule struct {
	CreatedAt			string   				`json:"created_at,omitempty"`
	UUID	    		string   				`json:"uuid,omitempty"`
	SourcePort 			int		 				`json:"source_port"`
	TargetPort			int		   				`json:"target_port"`
	Protocol			string 	 				`json:"protocol,omitempty"`
	Setting				ForwardingRuleSetting 	`json:"settings"`
}

type ForwardingRuleSetting struct {
	ConnectionLimit		int 	`json:"connection_limit"`
	SessionPersistence  string 	`json:"session_persistence"`
}

type ForwardingTarget struct {
	CreatedAt			string   `json:"created_at,omitempty"`
	TargetUUID	    	string   `json:"target_uuid"`
	TargetType			string   `json:"target_type"`
	TargetIPAddress 	string	 `json:"target_ip_address,omitempty"`
}

func (lb *LoadBalancerAPI) Init(c HTTPClient, authToken string, location string) error {
	lb.c = c
	lb.AuthToken = authToken
	lb.Location = location
	lb.ApiEndpoint = fmt.Sprintf(
		"https://api.idcloudhost.com/v1/%s/network/load_balancers",
		lb.Location,
	)
	r, err := http.Get(lb.ApiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("location: %s not found", lb.Location)
	}
	return nil
}

func (lb *LoadBalancerAPI) ListAll(isAll bool) error {
	url := fmt.Sprintf("%s?all=%t", lb.ApiEndpoint, isAll)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", lb.AuthToken)
	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	defer r.Body.Close()
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	bodyByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	
	return json.Unmarshal(bodyByte, &lb.LoadBalancerList)
}