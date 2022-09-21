package idcloudhost

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"bytes"
)

type LoadBalancerAPI struct {
	c           		HTTPClient
	AuthToken   		string
	Location    		string
	ApiEndpoint 		string
	LoadBalancer		*LoadBalancer
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
	ReservePublicIP 	bool				`json:"reserve_public_ip,omitempty"`
	ForwardingRules 	[]ForwardingRule 	`json:"forwarding_rules"`
	Targets				[]ForwardingTarget 	`json:"targets"`
}

type ForwardingRule struct {
	CreatedAt			string   				`json:"created_at,omitempty"`
	UUID	    		string   				`json:"uuid,omitempty"`
	SourcePort 			int		 				`json:"source_port"`
	TargetPort			int		   				`json:"target_port"`
	Protocol			string 	 				`json:"protocol,omitempty"`
	Setting				ForwardingRuleSetting 	`json:"settings,omitempty"`
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

func (lb *LoadBalancerAPI) Create(isAll bool, newLB *LoadBalancer) error {
	url := fmt.Sprintf("%s?all=%t", lb.ApiEndpoint, isAll)
	newLbJSON, err := json.Marshal(newLB)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newLbJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", lb.AuthToken)
	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	if err = json.NewDecoder(r.Body).Decode(&lb.LoadBalancer); err != nil {
		return err
	}
	return nil
}

func (lb *LoadBalancerAPI) Delete(UUID string) error {
	url := fmt.Sprintf("%s/%s", lb.ApiEndpoint, UUID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", lb.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	
	return checkError(r.StatusCode)
}

func (lb *LoadBalancerAPI) AddForwardingTarget(
		LBUUID string, TargetUUID string, TargetType string) error {
	url := fmt.Sprintf("%s/%s/targets", lb.ApiEndpoint, LBUUID)
	target := ForwardingTarget{
		TargetUUID: TargetUUID,
		TargetType: TargetType,
	}
	targetJSON, err := json.Marshal(&target)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(targetJSON))
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", lb.AuthToken)
	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	if err = checkError(r.StatusCode); err != nil {
		return err
	}
	if err = json.NewDecoder(r.Body).Decode(&target); err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

func (lb *LoadBalancerAPI) UnlinkForwardingTarget(LBUUID string, TargetUUID string) error {
	url := fmt.Sprintf("%s/%s/targets/%s", lb.ApiEndpoint, LBUUID, TargetUUID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("apiKey", lb.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	r, err := lb.c.Do(req)
	if err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	
	return checkError(r.StatusCode)
}


