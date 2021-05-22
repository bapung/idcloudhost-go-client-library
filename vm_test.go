package main

import (
	"log"
	"testing"
)

const userAuthToken = "wkkw"

func TestGetVMbyUUID(t *testing.T) {
	targetUuid := "wjwjw"
	loc := "jkt01"
	v := VirtualMachineAPI{}
	v.Init(userAuthToken, loc)
	if err := v.Get(targetUuid); err != nil {
		t.Fatal(err)
	}
	log.Println(v.VM)
}
