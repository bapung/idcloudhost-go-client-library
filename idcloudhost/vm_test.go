package idcloudhost

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

const userAuthToken = ""

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
