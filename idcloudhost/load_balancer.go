package idcloudhost

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoadBalancerAPI struct {
	c           		HTTPClient
	AuthToken   		string
	Location    		string
	ApiEndpoint 		string
	LoadBalancerList    *[]LoadBalancer
	LoadBalancer        *LoadBalancer
}


type LoadBalancer struct  {
	UUID 				string 				`json:"uuid"`
	NetworkUUID 		string 				`json:"network_uuid"`
	DisplayName			string				`json:"display_name"`
	UserUID 			int 				`json:"user_id"`
	BillingAccount 		int 				`json:"billing_account_id"`
	CreatedAt			string  			`json:"created_at"`
	UpdatedAt			string  			`json:"updated_at"`
	IsDeleted			bool    			`json:"is_deleted"`
	PrivateIPv4     	string  			`json:"private_address"`
	ReservePublicIP 	string				`json:"reserve_public_ip,omitempty"`
	ForwardingRules 	*[]ForwardingRule 	`json:"forwarding_rules"`
	Targets				*[]ForwardingTarget `json:"targets"`
}

type ForwardingRule struct {
	CreatedAt			string   				`json:"created_at,omitempty"`
	UUID	    		string   				`json:"uuid,omitempty"`
	SourcePort 			string	 				`json:"source_port"`
	TargetPort			string   				`json:"target_port"`
	Protocol			string 	 				`json:"protocol,omitempty"`
	Setting				*ForwardingRuleSetting 	`json:"settings"`
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