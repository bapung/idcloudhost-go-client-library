package main

import (
	"github.com/bapung/idcloudhost-go-client-library/idcloudhost"
)

func main() {
	u := idcloudhost.UserAPI{}
	u.Init("test")
}
