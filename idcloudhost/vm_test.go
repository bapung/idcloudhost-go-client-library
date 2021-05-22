package idcloudhost

import (
	"log"
	"testing"
)

const userAuthToken = "jUeD8GC6bx3esE8LjGQutPEZYnPMjNxa"

func TestGetVMbyUUID(t *testing.T) {
	targetUuid := "a28b4e97-c648-44ed-8217-f9d066dc6a91"
	loc := "jkt01"
	v := VirtualMachineAPI{}
	v.Init(userAuthToken, loc)
	if err := v.Get(targetUuid); err != nil {
		t.Fatal(err)
	}
	log.Println(v.VM)
}

func TestListAllVMs(t *testing.T) {
	loc := "jkt01"
	v := VirtualMachineAPI{}
	v.Init(userAuthToken, loc)
	if err := v.ListAll(); err != nil {
		t.Fatal(err)
	}
	log.Println(v.VMList[0])
}
