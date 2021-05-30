package idcloudhost

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

const userAuthToken = "RSwi3z5koSZ1hi7qA8zm50SswKDXJSN5"

func TestGetVMbyUUID(t *testing.T) {
	targetUuid := "a28b4e97-c648-44ed-8217-f9d066dc6a91"
	loc := "jkt01"
	v := VirtualMachineAPI{}
	v.Init(userAuthToken, loc)
	if err := v.Get(targetUuid); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(v.VMMap)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(b))
}

func TestListAllVMs(t *testing.T) {
	loc := "jkt01"
	v := VirtualMachineAPI{}
	v.Init(userAuthToken, loc)
	if err := v.ListAll(); err != nil {
		t.Fatal(err)
	}
	var inInterface []map[string]interface{}
	inrec, _ := json.Marshal(v.VMListMap)
	log.Println(json.Unmarshal(inrec, &inInterface))

	// iterate through inrecs
	for field, val := range inInterface {
		log.Println("KV Pair: ", field, val)
	}
}
func TestCreateVM(t *testing.T) {
	loc := "jkt01"
	v := VirtualMachineAPI{}
	v.Init(userAuthToken, loc)
	newVM := map[string]interface{}{
		"backup":          false,
		"name":            "testvm",
		"os_name":         "ubuntu",
		"os_version":      "16.04",
		"disks":           20,
		"vcpu":            1,
		"ram":             1024,
		"username":        "example",
		"password":        "Password123",
		"billing_account": 1200132376,
	}
	if err := v.Create(newVM); err != nil {
		t.Fatal(err)
	}
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(v.VMMap)
	log.Println(json.Unmarshal(inrec, &inInterface))

	// iterate through inrecs
	for field, val := range inInterface {
		log.Println("KV Pair: ", field, val)
	}
}
