package srvdiscovery

import (
	"fmt"
	"testing"
	"time"
)

func Test_discovery(t *testing.T) {
	Add("user", "/user/")
	Add("msg", "/msg/")
	if err := Watch(`{"endpoints":["localhost:2379"]}`); err != nil {
		t.Fatalf(err.Error())
	}

	for {
		fmt.Println("---->addr: ", Get("user"))
		fmt.Println("---->addr: ", Get("msg"))
		time.Sleep(1 * time.Second)
	}
}
